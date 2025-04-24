func MarshallToResponse(
    ctx context.Context,
    logger *slog.Logger,
    w http.ResponseWriter,
    v interface{}
){
	b, err := json.Marshal(v)
	if err != nil {
		logger.ErrorContext(
            ctx,
            "failed to serialize entity to json",
            slog.Any("error", err)
        )
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(b)
}
