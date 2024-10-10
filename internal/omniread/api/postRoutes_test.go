package api

// func TestDeletePostKnown(t *testing.T) {
// 	mockedRepo := &mockPostRepo{
// 		deleteFunc: func(ctx context.Context, id snowflake.Snowflake) error {
// 			return nil
// 		},
// 	}
//
// 	req, err := http.NewRequest("DELETE", "/post/1796290045997481984", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		mockedRepo,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 	}
// }
//
// func TestDeletePostUnknown(t *testing.T) {
// 	mockedRepo := &mockPostRepo{
// 		deleteFunc: func(ctx context.Context, id snowflake.Snowflake) error {
// 			return storage.NewNotFoundError(storage.Post, snowflake.ParseId(1796290045997481984))
// 		},
// 	}
//
// 	req, err := http.NewRequest("DELETE", "/post/1796290045997481984", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		mockedRepo,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusNotFound {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusNotFound)
// 	}
// }
//
// func TestDeletePostBadFormedId(t *testing.T) {
// 	req, err := http.NewRequest("DELETE", "/post/hello", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		nil,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusBadRequest)
// 	}
//
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("handler returned wrong content type: got %v want %v",
// 			rr.Header().Get("Content-Type"), "application/json")
// 	}
//
// 	expected := `{"error":"Bad Request","message":"Url parameter could not be parsed properly."}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestDeletePostDBError(t *testing.T) {
// 	mockedRepo := &mockPostRepo{
// 		deleteFunc: func(ctx context.Context, id snowflake.Snowflake) error {
// 			return storage.NewDatabaseError("database error", errors.New("database error"))
// 		},
// 	}
//
// 	req, err := http.NewRequest("DELETE", "/post/1796290045997481984", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		mockedRepo,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusInternalServerError {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusInternalServerError)
// 	}
// }
//
// func TestCreatePostSuccess(t *testing.T) {
// 	mockedRepo := &mockPostRepo{
// 		createFunc: func(ctx context.Context, entity models.Post) error {
// 			return nil
// 		},
// 	}
//
// 	body := struct {
// 		Id          uint64 `json:"id"`
// 		AuthorId    uint64 `json:"authorId"`
// 		AuthorName  string `json:"authorName"`
// 		Timestamp   string `json:"timestamp"`
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		ContentFile string `json:"contentFileUrl"`
// 	}{
// 		Id:          1796290045997481984,
// 		AuthorId:    1796290045997481985,
// 		AuthorName:  "johndoe",
// 		Timestamp:   "2021-01-01T11:40:35Z",
// 		Title:       "Hello, World!",
// 		Description: "Foobarbaz",
// 		ContentFile: "https://example.com/foo",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug})),
// 		nil,
// 		mockedRepo,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusCreated)
// 	}
//
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("handler returned wrong content type: got %v want %v",
// 			rr.Header().Get("Content-Type"), "application/json")
// 	}
//
// 	expected := `{"id":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2021-01-01T11:40:35Z","title":"Hello, World!","description":"Foobarbaz","contentFileUrl":"https://example.com/foo","comments":[],"tags":[]}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestCreatePostDuplicate(t *testing.T) {
// 	mockedRepo := &mockPostRepo{
// 		createFunc: func(ctx context.Context, entity models.Post) error {
// 			return storage.NewEntityAlreadyExistsError(snowflake.ParseId(1796290045997481984))
// 		},
// 	}
//
// 	body := struct {
// 		Id          uint64 `json:"id"`
// 		AuthorId    uint64 `json:"authorId"`
// 		AuthorName  string `json:"authorName"`
// 		Timestamp   string `json:"timestamp"`
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		ContentFile string `json:"contentFileUrl"`
// 	}{
// 		Id:          1796290045997481984,
// 		AuthorId:    1796290045997481985,
// 		AuthorName:  "johndoe",
// 		Timestamp:   "2021-01-01T11:40:35Z",
// 		Title:       "Hello, World!",
// 		Description: "Foobarbaz",
// 		ContentFile: "https://example.com/foo",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		mockedRepo,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusConflict {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusConflict)
// 	}
//
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("handler returned wrong content type: got %v want %v",
// 			rr.Header().Get("Content-Type"), "application/json")
// 	}
//
// 	expected := `{"error":"Conflict","message":"Post with that ID already exists."}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestCreatePostBadFormedBody(t *testing.T) {
// 	body := struct {
// 		Foo uint64 `json:"foo"`
// 		Bar string `json:"bar"`
// 		Baz string `json:"baz"`
// 	}{
// 		Foo: 1796290045997481984,
// 		Bar: "johndoe",
// 		Baz: "foobar",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		nil,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusBadRequest)
// 	}
//
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("handler returned wrong content type: got %v want %v",
// 			rr.Header().Get("Content-Type"), "application/json")
// 	}
//
// 	expected := `{"error":"Bad Request","message":"Request body could not be parsed properly."}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestCreatePostBadFormedTimestamp(t *testing.T) {
// 	body := struct {
// 		Id          uint64 `json:"id"`
// 		AuthorId    uint64 `json:"authorId"`
// 		AuthorName  string `json:"authorName"`
// 		Timestamp   string `json:"timestamp"`
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		ContentFile string `json:"contentFileUrl"`
// 	}{
// 		Id:          1796290045997481984,
// 		AuthorId:    1796290045997481985,
// 		AuthorName:  "johndoe",
// 		Timestamp:   "hello",
// 		Title:       "Hello, World!",
// 		Description: "Foobarbaz",
// 		ContentFile: "https://example.com/foo",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		nil,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusBadRequest)
// 	}
//
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("handler returned wrong content type: got %v want %v",
// 			rr.Header().Get("Content-Type"), "application/json")
// 	}
//
// 	expected := `{"error":"Bad Request","message":"Timestamp could not be parsed properly."}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestCreatePostDBError(t *testing.T) {
// 	mockedRepo := &mockPostRepo{
// 		createFunc: func(ctx context.Context, entity models.Post) error {
// 			return storage.NewDatabaseError("database error", errors.New("database error"))
// 		},
// 	}
//
// 	body := struct {
// 		Id          uint64 `json:"id"`
// 		AuthorId    uint64 `json:"authorId"`
// 		AuthorName  string `json:"authorName"`
// 		Timestamp   string `json:"timestamp"`
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		ContentFile string `json:"contentFileUrl"`
// 	}{
// 		Id:          1796290045997481984,
// 		AuthorId:    1796290045997481985,
// 		AuthorName:  "johndoe",
// 		Timestamp:   "2021-01-01T11:40:35Z",
// 		Title:       "Hello, World!",
// 		Description: "Foobarbaz",
// 		ContentFile: "https://example.com/foo",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", "/post", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		mockedRepo,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusInternalServerError {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusInternalServerError)
// 	}
// }
//
// func TestUpdatePostSuccess(t *testing.T) {
// 	mockedRepo := &mockPostRepo{
// 		updateFunc: func(ctx context.Context, entity models.Post) error {
// 			return nil
// 		},
// 	}
//
// 	body := struct {
// 		Id          uint64 `json:"id"`
// 		AuthorId    uint64 `json:"authorId"`
// 		AuthorName  string `json:"authorName"`
// 		Timestamp   string `json:"timestamp"`
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		ContentFile string `json:"contentFileUrl"`
// 	}{
// 		Id:          1796290045997481984,
// 		AuthorId:    1796290045997481985,
// 		AuthorName:  "johndoe",
// 		Timestamp:   "2021-01-01T11:40:35Z",
// 		Title:       "Hello, World!",
// 		Description: "Foobarbaz",
// 		ContentFile: "https://example.com/foo",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		mockedRepo,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusCreated {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusCreated)
// 	}
//
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("handler returned wrong content type: got %v want %v",
// 			rr.Header().Get("Content-Type"), "application/json")
// 	}
//
// 	expected := `{"id":1796290045997481984,"authorId":1796290045997481985,"authorName":"johndoe","timestamp":"2021-01-01T11:40:35Z","title":"Hello, World!","description":"Foobarbaz","contentFileUrl":"https://example.com/foo","comments":[],"tags":[]}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestUpdatePostNotFound(t *testing.T) {
// 	mockedRepo := &mockPostRepo{
// 		updateFunc: func(ctx context.Context, entity models.Post) error {
// 			return storage.NewNotFoundError(storage.Post, entity.Id())
// 		},
// 	}
//
// 	body := struct {
// 		Id          uint64 `json:"id"`
// 		AuthorId    uint64 `json:"authorId"`
// 		AuthorName  string `json:"authorName"`
// 		Timestamp   string `json:"timestamp"`
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		ContentFile string `json:"contentFileUrl"`
// 	}{
// 		Id:          1796290045997481984,
// 		AuthorId:    1796290045997481985,
// 		AuthorName:  "johndoe",
// 		Timestamp:   "2021-01-01T11:40:35Z",
// 		Title:       "Hello, World!",
// 		Description: "Foobarbaz",
// 		ContentFile: "https://example.com/foo",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		mockedRepo,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusNotFound {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusNotFound)
// 	}
//
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("handler returned wrong content type: got %v want %v",
// 			rr.Header().Get("Content-Type"), "application/json")
// 	}
//
// 	expected := `{"error":"Not Found","message":"Post with that ID could not be found to update."}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestUpdatePostBadFormedBody(t *testing.T) {
// 	body := struct {
// 		Foo uint64 `json:"foo"`
// 		Bar string `json:"bar"`
// 		Baz string `json:"baz"`
// 	}{
// 		Foo: 1796290045997481984,
// 		Bar: "johndoe",
// 		Baz: "foobar",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		nil,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusBadRequest)
// 	}
//
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("handler returned wrong content type: got %v want %v",
// 			rr.Header().Get("Content-Type"), "application/json")
// 	}
//
// 	expected := `{"error":"Bad Request","message":"Request body could not be parsed properly."}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestUpdatePostBadFormedTimestamp(t *testing.T) {
// 	body := struct {
// 		Id          uint64 `json:"id"`
// 		AuthorId    uint64 `json:"authorId"`
// 		AuthorName  string `json:"authorName"`
// 		Timestamp   string `json:"timestamp"`
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		ContentFile string `json:"contentFileUrl"`
// 	}{
// 		Id:          1796290045997481984,
// 		AuthorId:    1796290045997481985,
// 		AuthorName:  "johndoe",
// 		Timestamp:   "hello",
// 		Title:       "Hello, World!",
// 		Description: "Foobarbaz",
// 		ContentFile: "https://example.com/foo",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		nil,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusBadRequest {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusBadRequest)
// 	}
//
// 	if rr.Header().Get("Content-Type") != "application/json" {
// 		t.Errorf("handler returned wrong content type: got %v want %v",
// 			rr.Header().Get("Content-Type"), "application/json")
// 	}
//
// 	expected := `{"error":"Bad Request","message":"Timestamp could not be parsed properly."}`
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %v want %v",
// 			rr.Body.String(), expected)
// 	}
// }
//
// func TestUpdatePostDBError(t *testing.T) {
// 	mockedRepo := &mockPostRepo{
// 		updateFunc: func(ctx context.Context, entity models.Post) error {
// 			return storage.NewDatabaseError("database error", errors.New("database error"))
// 		},
// 	}
//
// 	body := struct {
// 		Id          uint64 `json:"id"`
// 		AuthorId    uint64 `json:"authorId"`
// 		AuthorName  string `json:"authorName"`
// 		Timestamp   string `json:"timestamp"`
// 		Title       string `json:"title"`
// 		Description string `json:"description"`
// 		ContentFile string `json:"contentFileUrl"`
// 	}{
// 		Id:          1796290045997481984,
// 		AuthorId:    1796290045997481985,
// 		AuthorName:  "johndoe",
// 		Timestamp:   "2021-01-01T11:40:35Z",
// 		Title:       "Hello, World!",
// 		Description: "Foobarbaz",
// 		ContentFile: "https://example.com/foo",
// 	}
// 	out, err := json.Marshal(body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	req, err := http.NewRequest("PUT", "/post/1796290045997481984", bytes.NewBuffer(out))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	rr := httptest.NewRecorder()
// 	handler := NewHandler(
// 		slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),
// 		nil,
// 		mockedRepo,
// 		nil,
// 	)
//
// 	handler.ServeHTTP(rr, req)
//
// 	if status := rr.Code; status != http.StatusInternalServerError {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusInternalServerError)
// 	}
// }
