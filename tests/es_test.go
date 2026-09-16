package tests

import (
	"context"
	"encoding/json"
	appes "gin/app/es"
	"gin/app/facade"
	"gin/app/model"
	"gin/common/ctxkey"
	"gin/config"
	eslib "gin/pkg/serviceprovider/es"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// fakeESStore 模拟ES存储
type fakeESStore struct {
	mu      sync.RWMutex
	indexes map[string]map[string]map[string]any
}

// newFakeESServer 创建模拟ES服务
func newFakeESServer() *httptest.Server {
	store := &fakeESStore{
		indexes: make(map[string]map[string]map[string]any),
	}
	return httptest.NewServer(store)
}

// ServeHTTP 处理模拟ES请求
func (s *fakeESStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if path == "" {
		writeESResponse(w, http.StatusOK, map[string]any{"name": "fake-es"})
		return
	}

	if len(parts) == 1 {
		s.handleIndex(w, r, parts[0])
		return
	}

	if len(parts) == 2 && parts[1] == "_search" {
		s.handleSearch(w, r, parts[0])
		return
	}

	if len(parts) >= 3 {
		switch parts[1] {
		case "_doc":
			s.handleDocument(w, r, parts[0], parts[2])
			return
		case "_update":
			s.handleUpdate(w, r, parts[0], parts[2])
			return
		case "_search":
			s.handleSearch(w, r, parts[0])
			return
		}
	}

	http.NotFound(w, r)
}

// handleIndex 处理索引操作
func (s *fakeESStore) handleIndex(w http.ResponseWriter, r *http.Request, index string) {
	switch r.Method {
	case http.MethodHead:
		s.mu.RLock()
		_, exists := s.indexes[index]
		s.mu.RUnlock()
		if exists {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	case http.MethodPut:
		s.mu.Lock()
		if _, exists := s.indexes[index]; !exists {
			s.indexes[index] = make(map[string]map[string]any)
		}
		s.mu.Unlock()
		writeESResponse(w, http.StatusOK, map[string]any{"acknowledged": true})
	case http.MethodDelete:
		s.mu.Lock()
		delete(s.indexes, index)
		s.mu.Unlock()
		writeESResponse(w, http.StatusOK, map[string]any{"acknowledged": true})
	default:
		http.NotFound(w, r)
	}
}

// handleDocument 处理文档索引和查询
func (s *fakeESStore) handleDocument(w http.ResponseWriter, r *http.Request, index, id string) {
	switch r.Method {
	case http.MethodPut:
		var doc map[string]any
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			writeESResponse(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}

		s.mu.Lock()
		if _, exists := s.indexes[index]; !exists {
			s.indexes[index] = make(map[string]map[string]any)
		}
		s.indexes[index][id] = doc
		s.mu.Unlock()
		writeESResponse(w, http.StatusOK, map[string]any{"result": "created"})
	case http.MethodGet:
		s.mu.RLock()
		doc, exists := s.indexes[index][id]
		s.mu.RUnlock()
		if !exists {
			writeESResponse(w, http.StatusNotFound, map[string]any{
				"_index": index,
				"_id":    id,
				"found":  false,
			})
			return
		}
		writeESResponse(w, http.StatusOK, map[string]any{
			"_index":  index,
			"_id":     id,
			"found":   true,
			"_source": doc,
		})
	case http.MethodDelete:
		s.mu.Lock()
		delete(s.indexes[index], id)
		s.mu.Unlock()
		writeESResponse(w, http.StatusOK, map[string]any{"result": "deleted"})
	default:
		http.NotFound(w, r)
	}
}

// handleUpdate 处理文档更新
func (s *fakeESStore) handleUpdate(w http.ResponseWriter, r *http.Request, index, id string) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	var body struct {
		Doc map[string]any `json:"doc"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeESResponse(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	s.mu.Lock()
	if _, exists := s.indexes[index]; !exists {
		s.indexes[index] = make(map[string]map[string]any)
	}
	doc, exists := s.indexes[index][id]
	if !exists {
		doc = make(map[string]any)
		s.indexes[index][id] = doc
	}
	maps.Copy(doc, body.Doc)
	s.mu.Unlock()
	writeESResponse(w, http.StatusOK, map[string]any{"result": "updated"})
}

// handleSearch 处理搜索
func (s *fakeESStore) handleSearch(w http.ResponseWriter, r *http.Request, index string) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	s.mu.RLock()
	hits := make([]map[string]any, 0, len(s.indexes[index]))
	for id, doc := range s.indexes[index] {
		hits = append(hits, map[string]any{
			"_index":  index,
			"_id":     id,
			"_score":  1,
			"_source": doc,
		})
	}
	s.mu.RUnlock()

	writeESResponse(w, http.StatusOK, map[string]any{
		"took":      1,
		"timed_out": false,
		"hits": map[string]any{
			"total": map[string]any{
				"value":    len(hits),
				"relation": "eq",
			},
			"hits": hits,
		},
	})
}

// writeESResponse 输出模拟ES响应
func writeESResponse(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// newTestESClient 创建测试ES客户端
func newTestESClient(t *testing.T) *eslib.Client {
	server := newFakeESServer()
	t.Cleanup(server.Close)

	client := eslib.NewClient(config.Es{
		Addresses: []string{server.URL},
		Timeout:   5 * time.Second,
	})
	return client
}

// registerTestESClient 注册测试ES客户端
func registerTestESClient(t *testing.T, client *eslib.Client) {
	previous := facade.ES()
	facade.Register[*eslib.Client]("es", client)
	t.Cleanup(func() {
		if previous != nil {
			facade.Register[*eslib.Client]("es", previous)
			return
		}
		facade.Unregister("es")
	})
}

// TestESClientCRUD 测试ES客户端CRUD
func TestESClientCRUD(t *testing.T) {
	client := newTestESClient(t)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-es-client-crud")

	index := "test_es_client_crud"
	defer func() {
		_ = client.DeleteIndex(ctx, index)
	}()

	require.NoError(t, client.Ping(ctx))

	exists, err := client.IndexExists(ctx, index)
	require.NoError(t, err)
	require.False(t, exists)

	require.NoError(t, client.CreateIndex(ctx, index, map[string]string{
		"username": "keyword",
		"nickname": "text",
	}))

	exists, err = client.IndexExists(ctx, index)
	require.NoError(t, err)
	require.True(t, exists)

	require.NoError(t, client.Index(ctx, index, "1", map[string]any{
		"username": "zhangsan",
		"nickname": "小张",
	}))

	doc, err := client.GetDocument[map[string]any](ctx, index, "1")
	require.NoError(t, err)
	require.True(t, doc.Found)
	require.Equal(t, "zhangsan", doc.Source["username"])

	require.NoError(t, client.Update(ctx, index, "1", map[string]any{
		"nickname": "张三",
	}))

	doc, err = client.GetDocument[map[string]any](ctx, index, "1")
	require.NoError(t, err)
	require.Equal(t, "张三", doc.Source["nickname"])

	result, err := client.Search[map[string]any](ctx, index, map[string]any{
		"query": map[string]any{
			"term": map[string]any{"username": "zhangsan"},
		},
	})
	require.NoError(t, err)
	require.Equal(t, int64(1), result.Hits.Total.Value)
	require.Len(t, result.Hits.Hits, 1)
	require.Equal(t, "zhangsan", result.Hits.Hits[0].Source["username"])

	require.NoError(t, client.Delete(ctx, index, "1"))
	_, err = client.GetDocument[map[string]any](ctx, index, "1")
	require.Error(t, err)

	require.NoError(t, client.DeleteIndex(ctx, index))
}

// TestESRequestTraceOnlyRecordsES 测试ES请求只记录ES调试信息
func TestESRequestTraceOnlyRecordsES(t *testing.T) {
	client := newTestESClient(t)
	traceID := "test-es-http-trace"
	store := facade.Debugger().Store()
	store.Delete(traceID)
	t.Cleanup(func() {
		store.Delete(traceID)
	})

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, traceID)

	require.NoError(t, client.Ping(ctx))

	trace, ok := store.Get(traceID)
	require.True(t, ok)
	require.NotEmpty(t, trace.ES)
	require.Empty(t, trace.HTTP)
}

// TestESUserSearch 测试用户ES搜索封装
func TestESUserSearch(t *testing.T) {
	client := newTestESClient(t)
	registerTestESClient(t, client)

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-es-user-search")

	search := &appes.UserSearch{}
	mapping := search.Mapping()
	require.Equal(t, "nested", mapping["userRoles"])
	require.Equal(t, "date", mapping["createdAt"])

	query := search.Query(map[string]any{"keyword": "zhangsan"})
	require.Contains(t, query, "multi_match")

	defer func() {
		_ = search.DeleteIndex(ctx)
	}()

	require.NoError(t, search.EnsureIndex(ctx))

	createdAt := model.DateTime(time.Date(2026, 9, 11, 10, 0, 0, 0, time.Local))
	user := &model.User{
		ID:       1,
		Username: "zhangsan",
		FullName: "张三",
		Nickname: "小张",
		Email:    "zhangsan@example.com",
		Gender:   1,
		Age:      18,
		Status:   1,
		UserRoles: []*model.UserRoles{
			{ID: 1, UserID: 1, RoleID: 1, Name: "管理员"},
		},
		CreatedAt: &createdAt,
		UpdatedAt: &createdAt,
	}

	require.NoError(t, search.Save(ctx, user))

	detail, err := search.Detail(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "zhangsan", detail["username"])
	require.Equal(t, "小张", detail["nickname"])
	require.NotEmpty(t, detail["createdAt"])

	require.NoError(t, search.Update(ctx, user.ID, map[string]any{
		"nickname": "张三",
		"age":      19,
	}))

	detail, err = search.Detail(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, "张三", detail["nickname"])
	require.Equal(t, float64(19), detail["age"])

	page, err := search.List(ctx, map[string]any{"username": "zhangsan"}, 1, 10, nil)
	require.NoError(t, err)
	require.Equal(t, int64(1), page.Total)
	require.Len(t, page.List, 1)

	require.NoError(t, search.Delete(ctx, user.ID))
	_, err = search.Detail(ctx, user.ID)
	require.Error(t, err)

	require.NoError(t, search.DeleteIndex(ctx))
}
