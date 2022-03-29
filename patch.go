package log

// Patch patches the existing settings with any option given.
// This is thread safe but does not propagates to child loggers.
func (l *Logger) Patch(options ...Option) {
	l.settingsMutex.Lock()
	defer l.settingsMutex.Unlock()

	updatedSettings := l.settings.copy()
	for _, option := range options {
		option(&updatedSettings)
	}

	writerMutexes := writersRegistry.RegisterWriters(updatedSettings.writers)

	l.settings = updatedSettings
	l.writersMutexesMutex.Lock()
	l.writersMutexes = writerMutexes
	l.writersMutexesMutex.Unlock()
}
