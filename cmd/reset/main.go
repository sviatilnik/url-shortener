package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// StructInfo содержит информацию о структуре для генерации метода Reset
type StructInfo struct {
	Name    string
	Package string
	File    string
	Fields  []FieldInfo
}

// FieldInfo содержит информацию о поле структуры
type FieldInfo struct {
	Name      string
	Type      string
	IsPointer bool
	IsSlice   bool
	IsMap     bool
	IsStruct  bool
	Tag       string
}

func main() {
	// Получаем корневую директорию проекта
	rootDir := "."
	if len(os.Args) > 1 {
		rootDir = os.Args[1]
	}

	// Сканируем все пакеты
	structs, err := scanPackages(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning packages: %v\n", err)
		os.Exit(1)
	}

	// Группируем структуры по пакетам
	packageStructs := make(map[string][]StructInfo)
	for _, s := range structs {
		packageStructs[s.Package] = append(packageStructs[s.Package], s)
	}

	// Генерируем файлы reset.gen.go для каждого пакета
	for pkg, structs := range packageStructs {
		if err := generateResetFile(pkg, structs); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating reset file for package %s: %v\n", pkg, err)
			os.Exit(1)
		}
		fmt.Printf("Generated reset.gen.go for package %s (%d structs)\n", pkg, len(structs))
	}
}

// scanPackages сканирует все пакеты в директории и находит структуры с комментарием // generate:reset
func scanPackages(rootDir string) ([]StructInfo, error) {
	var structs []StructInfo

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем директории и файлы, которые не являются .go файлами
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Пропускаем уже сгенерированные файлы
		if strings.HasSuffix(path, "reset.gen.go") {
			return nil
		}

		// Пропускаем файлы в директориях vendor, .git, .vscode, .idea
		dir := filepath.Dir(path)
		if strings.Contains(dir, "vendor") || strings.Contains(dir, ".git") ||
			strings.Contains(dir, ".vscode") || strings.Contains(dir, ".idea") {
			return nil
		}

		// Парсим файл
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil // Пропускаем файлы с ошибками парсинга
		}

		// Получаем имя пакета
		packageName := node.Name.Name

		// Ищем структуры с комментарием // generate:reset
		ast.Inspect(node, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.GenDecl:
				if x.Tok == token.TYPE {
					for _, spec := range x.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							if structType, ok := typeSpec.Type.(*ast.StructType); ok {
								// Проверяем комментарии
								if hasGenerateResetComment(x.Doc) {
									fields := extractFields(structType)
									structs = append(structs, StructInfo{
										Name:    typeSpec.Name.Name,
										Package: packageName,
										File:    path,
										Fields:  fields,
									})
								}
							}
						}
					}
				}
			}
			return true
		})

		return nil
	})

	return structs, err
}

// hasGenerateResetComment проверяет, есть ли в комментариях // generate:reset
func hasGenerateResetComment(doc *ast.CommentGroup) bool {
	if doc == nil {
		return false
	}

	for _, comment := range doc.List {
		if strings.Contains(comment.Text, "// generate:reset") {
			return true
		}
	}
	return false
}

// extractFields извлекает информацию о полях структуры
func extractFields(structType *ast.StructType) []FieldInfo {
	var fields []FieldInfo

	for _, field := range structType.Fields.List {
		fieldType := getTypeString(field.Type)
		fieldName := ""
		if len(field.Names) > 0 {
			fieldName = field.Names[0].Name
		}

		fields = append(fields, FieldInfo{
			Name:      fieldName,
			Type:      fieldType,
			IsPointer: isPointer(field.Type),
			IsSlice:   isSlice(field.Type),
			IsMap:     isMap(field.Type),
			IsStruct:  isStruct(field.Type),
			Tag:       getTag(field.Tag),
		})
	}

	return fields
}

// getTypeString возвращает строковое представление типа
func getTypeString(expr ast.Expr) string {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name
	case *ast.StarExpr:
		return "*" + getTypeString(x.X)
	case *ast.ArrayType:
		if x.Len == nil {
			return "[]" + getTypeString(x.Elt)
		}
		return fmt.Sprintf("[%s]%s", getTypeString(x.Len), getTypeString(x.Elt))
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", getTypeString(x.Key), getTypeString(x.Value))
	case *ast.SelectorExpr:
		return getTypeString(x.X) + "." + x.Sel.Name
	default:
		return "unknown"
	}
}

func isPointer(expr ast.Expr) bool {
	_, ok := expr.(*ast.StarExpr)
	return ok
}

func isSlice(expr ast.Expr) bool {
	if arrayType, ok := expr.(*ast.ArrayType); ok {
		return arrayType.Len == nil
	}
	return false
}

func isMap(expr ast.Expr) bool {
	_, ok := expr.(*ast.MapType)
	return ok
}

func isStruct(expr ast.Expr) bool {
	// Убираем указатель если есть
	if starExpr, ok := expr.(*ast.StarExpr); ok {
		expr = starExpr.X
	}

	if ident, ok := expr.(*ast.Ident); ok {
		// Простая проверка - если это не встроенный тип, считаем структурой
		builtinTypes := map[string]bool{
			"int": true, "int8": true, "int16": true, "int32": true, "int64": true,
			"uint": true, "uint8": true, "uint16": true, "uint32": true, "uint64": true,
			"float32": true, "float64": true, "complex64": true, "complex128": true,
			"bool": true, "string": true, "byte": true, "rune": true,
			"interface": true, "error": true,
		}
		return !builtinTypes[ident.Name]
	}
	return false
}

// getTag возвращает тег поля
func getTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}
	return tag.Value
}

// generateResetFile генерирует файл reset.gen.go для пакета
func generateResetFile(packageName string, structs []StructInfo) error {
	if len(structs) == 0 {
		return nil
	}

	// Определяем директорию пакета
	packageDir := filepath.Dir(structs[0].File)
	outputFile := filepath.Join(packageDir, "reset.gen.go")

	// Генерируем содержимое файла
	var buf bytes.Buffer

	// Заголовок файла
	buf.WriteString("// Code generated by reset tool. DO NOT EDIT.\n\n")
	buf.WriteString(fmt.Sprintf("package %s\n\n", packageName))

	// Генерируем методы Reset для каждой структуры
	for _, s := range structs {
		buf.WriteString(generateResetMethod(s))
		buf.WriteString("\n")
	}

	// Форматируем код
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("error formatting generated code: %v", err)
	}

	// Записываем в файл
	return os.WriteFile(outputFile, formatted, 0644)
}

// generateResetMethod генерирует метод Reset для структуры
func generateResetMethod(s StructInfo) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("func (s *%s) Reset() {\n", s.Name))

	for _, field := range s.Fields {
		if field.Name == "" {
			continue // Пропускаем анонимные поля
		}

		if field.IsPointer {
			buf.WriteString(fmt.Sprintf("\tif s.%s != nil {\n", field.Name))
			// Для указателей на структуры проверяем, есть ли метод Reset
			if field.IsStruct {
				buf.WriteString(fmt.Sprintf("\t\tif resetter, ok := interface{}(s.%s).(interface{ Reset() }); ok {\n", field.Name))
				buf.WriteString(fmt.Sprintf("\t\t\tresetter.Reset()\n"))
				buf.WriteString(fmt.Sprintf("\t\t}\n"))
				buf.WriteString(fmt.Sprintf("\t\ts.%s = nil\n", field.Name))
			} else {
				buf.WriteString(generateFieldReset(field, "s."+field.Name, 2))
			}
			buf.WriteString("\t}\n")
		} else {
			buf.WriteString(generateFieldReset(field, "s."+field.Name, 1))
		}
	}

	buf.WriteString("}\n")
	return buf.String()
}

// generateFieldReset генерирует код сброса для конкретного поля
func generateFieldReset(field FieldInfo, fieldPath string, indentLevel int) string {
	indent := strings.Repeat("\t", indentLevel)
	var buf bytes.Buffer

	if field.IsSlice {
		// Слайсы обрезаем по длине, но не зануляем
		buf.WriteString(fmt.Sprintf("%s%s = %s[:0]\n", indent, fieldPath, fieldPath))
	} else if field.IsMap {
		// Мапы очищаем
		buf.WriteString(fmt.Sprintf("%sclear(%s)\n", indent, fieldPath))
	} else if field.IsStruct {
		// Для структур вызываем метод Reset, если он есть
		buf.WriteString(fmt.Sprintf("%s// Вызов Reset для структуры %s\n", indent, field.Type))
		buf.WriteString(fmt.Sprintf("%sif resetter, ok := interface{}(&%s).(interface{ Reset() }); ok {\n", indent, fieldPath))
		buf.WriteString(fmt.Sprintf("%s    resetter.Reset()\n", indent))
		buf.WriteString(fmt.Sprintf("%s} else {\n", indent))
		zeroValue := getZeroValue(field.Type)
		buf.WriteString(fmt.Sprintf("%s    %s = %s\n", indent, fieldPath, zeroValue))
		buf.WriteString(fmt.Sprintf("%s}\n", indent))
	} else {
		// Примитивные типы сбрасываем к нулевым значениям
		zeroValue := getZeroValue(field.Type)
		buf.WriteString(fmt.Sprintf("%s%s = %s\n", indent, fieldPath, zeroValue))
	}

	return buf.String()
}

// getZeroValue возвращает нулевое значение для типа
func getZeroValue(typeStr string) string {
	// Проверяем, является ли тип указателем
	isPointer := strings.HasPrefix(typeStr, "*")
	baseType := strings.TrimPrefix(typeStr, "*")

	var zeroValue string
	switch baseType {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"float32", "float64", "complex64", "complex128":
		zeroValue = "0"
	case "bool":
		zeroValue = "false"
	case "string":
		zeroValue = `""`
	case "byte", "rune":
		zeroValue = "0"
	default:
		// Для пользовательских типов возвращаем нулевое значение
		zeroValue = baseType + "{}"
	}

	// Если это указатель, добавляем & или nil
	if isPointer {
		return "nil"
	}
	return zeroValue
}
