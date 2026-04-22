package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

// Выталкивание элементов при переполнении.
func TestCacheEviction(t *testing.T) {
	// Создаём кэш вместимостью 3 элемента
	cache := NewCache(3)

	// Добавляем 3 элемента (кэш заполнен)
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	// Проверяем, что все 3 элемента на месте
	t.Run("cache filled with 3 items", func(t *testing.T) {
		val, ok := cache.Get("a")
		require.True(t, ok, "key 'a' should exist")
		require.Equal(t, 1, val)

		val, ok = cache.Get("b")
		require.True(t, ok, "key 'b' should exist")
		require.Equal(t, 2, val)

		val, ok = cache.Get("c")
		require.True(t, ok, "key 'c' should exist")
		require.Equal(t, 3, val)
	})

	// Добавляем 4-й элемент - должен вытолкнуться самый старый (a)
	cache.Set("d", 4)

	t.Run("adding 4th item should evict the oldest one (a)", func(t *testing.T) {
		// Проверяем, что 'a' был вытолкнут
		_, ok := cache.Get("a")
		require.False(t, ok, "key 'a' should be evicted")

		// Проверяем, что остальные элементы на месте
		val, ok := cache.Get("b")
		require.True(t, ok, "key 'b' should still exist")
		require.Equal(t, 2, val)

		val, ok = cache.Get("c")
		require.True(t, ok, "key 'c' should still exist")
		require.Equal(t, 3, val)

		val, ok = cache.Get("d")
		require.True(t, ok, "key 'd' should exist")
		require.Equal(t, 4, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

// Правильность LRU-порядка.
func TestCacheLRUOrder(t *testing.T) {
	cache := NewCache(3)

	// Добавляем элементы
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	// Используем элемент 'a' (он становится самым свежим)
	cache.Get("a")

	// Добавляем новый элемент - должен вытолкнуться самый старый (b)
	cache.Set("d", 4)

	t.Run("accessing 'a' makes it recent, 'b' becomes oldest and gets evicted", func(t *testing.T) {
		// 'b' должен быть вытолкнут (стал самым старым после использования 'a')
		_, ok := cache.Get("b")
		require.False(t, ok, "key 'b' should be evicted as it was the least recently used")

		// 'a' должен быть на месте (только что использовали)
		val, ok := cache.Get("a")
		require.True(t, ok, "key 'a' should exist (was recently used)")
		require.Equal(t, 1, val)

		// 'c' и 'd' должны быть на месте
		val, ok = cache.Get("c")
		require.True(t, ok, "key 'c' should exist")
		require.Equal(t, 3, val)

		val, ok = cache.Get("d")
		require.True(t, ok, "key 'd' should exist")
		require.Equal(t, 4, val)
	})
}

// Обновление существующего элемента.
func TestCacheUpdate(t *testing.T) {
	cache := NewCache(3)

	// Добавляем элементы
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	// Обновляем элемент 'b'
	old := cache.Set("b", 20)

	t.Run("updating existing key should move it to front", func(t *testing.T) {
		require.True(t, old, "Set should return true when updating existing key")

		// Проверяем обновлённое значение
		val, ok := cache.Get("b")
		require.True(t, ok)
		require.Equal(t, 20, val)
	})

	// Добавляем новый элемент - должен вытолкнуться самый старый (a)
	// После обновления 'b' стал свежим, самый старый - 'a'
	cache.Set("d", 4)

	t.Run("after update, 'b' is recent, so 'a' should be evicted", func(t *testing.T) {
		_, ok := cache.Get("a")
		require.False(t, ok, "key 'a' should be evicted (was never used after initial add)")

		// Остальные элементы должны быть на месте
		_, ok = cache.Get("b")
		require.True(t, ok, "key 'b' should exist (was just updated)")

		_, ok = cache.Get("c")
		require.True(t, ok, "key 'c' should exist")

		_, ok = cache.Get("d")
		require.True(t, ok, "key 'd' should exist")
	})
}

// Очистка кэша.
func TestCacheClear(t *testing.T) {
	cache := NewCache(3)

	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)

	cache.Clear()

	t.Run("after clear, cache should be empty", func(t *testing.T) {
		keys := []string{"a", "b", "c"}
		for _, key := range keys {
			_, ok := cache.Get(Key(key))
			require.False(t, ok, "key '%s' should not exist after Clear()", key)
		}
	})

	t.Run("after clear, cache should accept new items", func(t *testing.T) {
		cache.Set("d", 4)
		val, ok := cache.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)
	})
}

// Разные типы значений.
func TestCacheWithDifferentTypes(t *testing.T) {
	cache := NewCache(3)

	cache.Set("int", 42)
	cache.Set("string", "hello")
	cache.Set("struct", struct{ Name string }{"John"})

	t.Run("should store and retrieve different types correctly", func(t *testing.T) {
		val, ok := cache.Get("int")
		require.True(t, ok)
		require.Equal(t, 42, val)

		val, ok = cache.Get("string")
		require.True(t, ok)
		require.Equal(t, "hello", val)

		val, ok = cache.Get("struct")
		require.True(t, ok)
		require.Equal(t, struct{ Name string }{"John"}, val)
	})
}
