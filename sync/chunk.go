package sync

import "fmt"

// chunk represents file chunks. Files can be split into pieces and downloaded
// with multiple connections, each connection fetches a part of a file.
type chunk struct {
	// Where the chunk starts
	offset int64

	// Length of chunk
	length int64
}

// String implements fmt.Stringer for chunk.
func (c chunk) String() string {
	return fmt.Sprintf("chunk{%v-%v}", c.offset, c.offset+c.length)
}

// calculateChunks splits a filesize into count parts, in which every chunk is
// either smaller than blocksize, or can be fully divisible. Chunks has to be
// divisible by blocksize since we maintain a bitfield for each file, and each
// bit represents the said constant.
func calculateChunks(state *State, segmentnum uint) []*chunk {
	count := int64(segmentnum)
	pieceLength := int64(state.BitfieldPieceLength)

	// calculate the chunks of a resumable file.
	if state.Bitfield.Count() != 0 {
		var chunks []*chunk
		var idx uint32
		for {
			start, ok := state.Bitfield.FirstClear(idx)
			if !ok {
				break
			}

			end, ok := state.Bitfield.FirstSet(start)
			if !ok {
				chunks = append(chunks, &chunk{
					offset: int64(start) * pieceLength,
					length: state.FileLength - int64(start)*pieceLength,
				})
				break
			}

			chunks = append(chunks, &chunk{
				offset: int64(start) * pieceLength,
				length: int64(end-start) * pieceLength,
			})

			idx = end
		}
		return chunks
	}

	// calculate the chunks of a fresh new file.

	filesize := state.FileLength
	// don't even consider smaller files
	if filesize <= pieceLength || count <= 1 {
		return []*chunk{{offset: 0, length: filesize}}
	}

	// how many blocks fit perfectly on a filesize
	blockCount := filesize / pieceLength
	// how many bytes are left out
	excessBytes := filesize % pieceLength

	// If there are no blocks available for the given blocksize, we're gonna
	// reduce the count to the max available block count.
	if blockCount < count {
		count = blockCount
	}

	blocksPerUnit := blockCount / count
	excessBlocks := blockCount % count

	var chunks []*chunk
	for i := int64(0); i < count; i++ {
		chunks = append(chunks, &chunk{
			offset: i * blocksPerUnit * pieceLength,
			length: blocksPerUnit * pieceLength,
		})
	}

	if excessBlocks > 0 {
		offset := count * blocksPerUnit * pieceLength
		length := excessBlocks * pieceLength
		chunks = append(chunks, &chunk{
			offset: offset,
			length: length,
		})
	}

	// append excess bytes to the last chunk
	if excessBytes > 0 {
		c := chunks[len(chunks)-1]
		c.length += excessBytes
	}

	return chunks
}
