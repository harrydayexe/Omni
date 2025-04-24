func MarshallToResponse(
    ctx context.Context, logger *slog.Logger,
    w http.ResponseWriter, v interface{}
){
	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
