package main

// HashResult represents a single hash computation result
type HashResult struct {
	Type    string `json:"type"`
	Value   string `json:"value"`
	Bits    int    `json:"bits"`
	Binbits int    `json:"binbits,omitempty"` // Only for colorhash
}

// SegmentHash represents a segment in crop-resistant hash
type SegmentHash struct {
	Value string `json:"value"`
	Bits  int    `json:"bits"`
}

// CropResistantResult represents crop-resistant hash result
type CropResistantResult struct {
	Type          string        `json:"type"`
	Segments      int           `json:"segments"`
	SegmentHashes []SegmentHash `json:"segment_hashes,omitempty"`
}

// OutputData represents complete output data
type OutputData struct {
	Image  string                `json:"image"`
	Hash   *HashResult           `json:"hash,omitempty"`   // For single hash
	Hashes map[string]HashResult `json:"hashes,omitempty"` // For -all flag
	Crop   *CropResistantResult  `json:"crop,omitempty"`   // For crop-resistant hash
	Error  string                `json:"error,omitempty"`
}
