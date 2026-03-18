package errs

import "errors"

var ErrLLMNilFinalMessage = errors.New("llm returned nil final message")