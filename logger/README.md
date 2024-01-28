# Logger Package

This package provides a logging wrapper for the default golang slog package with some extra functionality, for example:
* Add color to the log level (blue for INFO, red for ERROR, etc...) 
* Add the `f` and `w` options for all methods, like `Infof` and `Infow`.
* A logParameter `AddAttributesFromContext(ctx) []any` to add default attributes from context.
