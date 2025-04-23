if errors.Is(err, sql.ErrNoRows) {
	logger.InfoContext(
        ctx, "entity not found", 
        slog.Any("id", id)
    )
	http.Error(
        w, "entity not found", 
        http.StatusNotFound
    )
	return true
}
