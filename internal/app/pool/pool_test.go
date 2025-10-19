package pool

import (
	"testing"
)

type TestStruct struct {
	ID      int
	Name    string
	Data    []string
	Counter int
}

func (t *TestStruct) Reset() {
	t.ID = 0
	t.Name = ""
	t.Data = t.Data[:0]
	t.Counter = 0
}

func TestPool_GetPut(t *testing.T) {
	// Создаем пул
	pool := New[*TestStruct](func() *TestStruct {
		return &TestStruct{}
	})

	// Получаем объект из пула
	obj1 := pool.Get()
	if obj1 == nil {
		t.Fatal("Expected non-nil object from pool")
	}

	// Заполняем объект данными
	obj1.ID = 123
	obj1.Name = "test"
	obj1.Data = []string{"item1", "item2"}
	obj1.Counter = 42

	// Возвращаем объект в пул
	pool.Put(obj1)

	// Получаем объект снова
	obj2 := pool.Get()
	if obj2 == nil {
		t.Fatal("Expected non-nil object from pool")
	}

	// Проверяем, что объект был сброшен
	if obj2.ID != 0 {
		t.Errorf("Expected ID 0, got %d", obj2.ID)
	}
	if obj2.Name != "" {
		t.Errorf("Expected empty Name, got '%s'", obj2.Name)
	}
	if len(obj2.Data) != 0 {
		t.Errorf("Expected empty Data slice, got length %d", len(obj2.Data))
	}
	if obj2.Counter != 0 {
		t.Errorf("Expected Counter 0, got %d", obj2.Counter)
	}
}

func TestPool_MultipleObjects(t *testing.T) {
	// Создаем пул
	pool := New[*TestStruct](func() *TestStruct {
		return &TestStruct{}
	})

	// Получаем несколько объектов
	obj1 := pool.Get()
	obj2 := pool.Get()

	// Проверяем, что это разные объекты
	if obj1 == obj2 {
		t.Error("Expected different objects from pool")
	}

	// Заполняем объекты разными данными
	obj1.ID = 1
	obj1.Name = "first"
	obj2.ID = 2
	obj2.Name = "second"

	// Возвращаем объекты в пул
	pool.Put(obj1)
	pool.Put(obj2)

	// Получаем объекты снова
	obj3 := pool.Get()
	obj4 := pool.Get()

	// Проверяем, что объекты были сброшены
	if obj3.ID != 0 || obj3.Name != "" {
		t.Error("Expected obj3 to be reset")
	}
	if obj4.ID != 0 || obj4.Name != "" {
		t.Error("Expected obj4 to be reset")
	}
}

func TestPool_Reuse(t *testing.T) {
	// Создаем пул
	pool := New[*TestStruct](func() *TestStruct {
		return &TestStruct{}
	})

	// Получаем объект
	obj1 := pool.Get()
	obj1.ID = 999
	obj1.Name = "original"

	// Возвращаем в пул
	pool.Put(obj1)

	// Получаем объект снова
	obj2 := pool.Get()

	// Проверяем, что это тот же объект (переиспользование)
	if obj1 != obj2 {
		t.Error("Expected same object to be reused")
	}

	// Проверяем, что объект был сброшен
	if obj2.ID != 0 || obj2.Name != "" {
		t.Error("Expected object to be reset before reuse")
	}
}

func TestPool_Concurrent(t *testing.T) {
	// Создаем пул
	pool := New[*TestStruct](func() *TestStruct {
		return &TestStruct{}
	})

	// Количество горутин
	numGoroutines := 100
	objects := make(chan *TestStruct, numGoroutines)

	// Запускаем горутины
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			// Получаем объект
			obj := pool.Get()
			obj.ID = id
			obj.Name = "goroutine"
			obj.Counter = id * 2

			// Возвращаем объект
			pool.Put(obj)
			objects <- obj
		}(i)
	}

	// Собираем результаты
	for i := 0; i < numGoroutines; i++ {
		obj := <-objects
		// Проверяем, что объект был сброшен
		if obj.ID != 0 || obj.Name != "" || obj.Counter != 0 {
			t.Errorf("Expected object to be reset, got ID=%d, Name='%s', Counter=%d",
				obj.ID, obj.Name, obj.Counter)
		}
	}
}

// BenchmarkPool_GetPut измеряет производительность операций Get/Put
func BenchmarkPool_GetPut(b *testing.B) {
	pool := New[*TestStruct](func() *TestStruct {
		return &TestStruct{}
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj := pool.Get()
		obj.ID = i
		obj.Name = "benchmark"
		pool.Put(obj)
	}
}

// BenchmarkPool_NewObject измеряет производительность создания новых объектов
func BenchmarkPool_NewObject(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		obj := &TestStruct{}
		obj.ID = i
		obj.Name = "benchmark"
		obj.Reset()
	}
}
