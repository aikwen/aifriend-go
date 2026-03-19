package rag

import "strings"

// SplitText 按固定长度切分文本，并保留相邻块之间的重叠部分。
// 例如长度是1400, chunkSize=500, overlap=50 时：
// 第1块 [0:500]
// 第2块 [450:950]
// 第3块 [900:1400]
func SplitText(text string, chunkSize int, overlap int) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}

	if chunkSize <= 0 {
		return []string{text}
	}

	if overlap < 0 {
		overlap = 0
	}

	// 防止 step <= 0 导致死循环
	if overlap >= chunkSize {
		overlap = chunkSize - 1
		if overlap < 0 {
			overlap = 0
		}
	}

	runes := []rune(text)
	if len(runes) <= chunkSize {
		return []string{text}
	}

	step := chunkSize - overlap
	chunks := make([]string, 0, (len(runes)+step-1)/step)

	for start := 0; start < len(runes); start += step {
		end := start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}

		chunk := strings.TrimSpace(string(runes[start:end]))
		if chunk != "" {
			chunks = append(chunks, chunk)
		}

		if end == len(runes) {
			break
		}
	}

	return chunks
}