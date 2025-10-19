package pool

import (
	"fmt"
	"testing"
)

func ExamplePool_usage() {
	userPool := New[*User](func() *User {
		return &User{}
	})

	// Получаем объект из пула
	user := userPool.Get()

	// Заполняем данными
	user.ID = 123
	user.Name = "John Doe"
	user.Email = "john@example.com"
	user.Active = true
	user.Tags = []string{"admin", "user"}
	user.Settings = map[string]string{
		"theme": "dark",
		"lang":  "en",
	}

	fmt.Printf("User before reset: ID=%d, Name=%s, Active=%v\n",
		user.ID, user.Name, user.Active)

	// Возвращаем объект в пул (автоматически вызовется Reset())
	userPool.Put(user)

	// Получаем объект снова
	reusedUser := userPool.Get()

	fmt.Printf("User after reset: ID=%d, Name=%s, Active=%v\n",
		reusedUser.ID, reusedUser.Name, reusedUser.Active)

	// Output:
	// User before reset: ID=123, Name=John Doe, Active=true
	// User after reset: ID=0, Name=, Active=false
}

// ExamplePool_multipleTypes демонстрирует использование Pool с разными типами
func ExamplePool_multipleTypes() {
	// Пул для User
	userPool := New[*User](func() *User {
		return &User{}
	})

	// Пул для Profile
	profilePool := New[*Profile](func() *Profile {
		return &Profile{}
	})

	// Получаем объекты из разных пулов
	user := userPool.Get()
	profile := profilePool.Get()

	// Заполняем данными
	user.ID = 1
	user.Name = "Alice"
	profile.Bio = "Software developer"
	profile.Avatar = "avatar.jpg"

	fmt.Printf("User: %s, Profile: %s\n", user.Name, profile.Bio)

	// Возвращаем в соответствующие пулы
	userPool.Put(user)
	profilePool.Put(profile)

	// Получаем объекты снова
	reusedUser := userPool.Get()
	reusedProfile := profilePool.Get()

	fmt.Printf("Reused User: %s, Reused Profile: %s\n",
		reusedUser.Name, reusedProfile.Bio)

	// Output:
	// User: Alice, Profile: Software developer
	// Reused User: , Reused Profile:
}

// TestPool_WithRealStructs тестирует Pool с реальными структурами из проекта
func TestPool_WithRealStructs(t *testing.T) {
	// Создаем пул для User
	userPool := New[*User](func() *User {
		return &User{}
	})

	// Получаем объект
	user := userPool.Get()
	if user == nil {
		t.Fatal("Expected non-nil user")
	}

	// Заполняем данными
	user.ID = 456
	user.Name = "Test User"
	user.Email = "test@example.com"
	user.Active = true
	user.Tags = []string{"test", "user"}
	user.Settings = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	// Проверяем, что данные установлены
	if user.ID != 456 {
		t.Errorf("Expected ID 456, got %d", user.ID)
	}
	if len(user.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(user.Tags))
	}
	if len(user.Settings) != 2 {
		t.Errorf("Expected 2 settings, got %d", len(user.Settings))
	}

	// Возвращаем в пул
	userPool.Put(user)

	// Получаем объект снова
	reusedUser := userPool.Get()
	if reusedUser == nil {
		t.Fatal("Expected non-nil reused user")
	}

	// Проверяем, что объект был сброшен
	if reusedUser.ID != 0 {
		t.Errorf("Expected ID 0, got %d", reusedUser.ID)
	}
	if reusedUser.Name != "" {
		t.Errorf("Expected empty Name, got '%s'", reusedUser.Name)
	}
	if reusedUser.Email != "" {
		t.Errorf("Expected empty Email, got '%s'", reusedUser.Email)
	}
	if reusedUser.Active != false {
		t.Errorf("Expected Active false, got %v", reusedUser.Active)
	}
	if len(reusedUser.Tags) != 0 {
		t.Errorf("Expected empty Tags slice, got length %d", len(reusedUser.Tags))
	}
	if len(reusedUser.Settings) != 0 {
		t.Errorf("Expected empty Settings map, got length %d", len(reusedUser.Settings))
	}
}
