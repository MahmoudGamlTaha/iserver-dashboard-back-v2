package repositories

import (
	"fmt"

	"github.com/google/uuid"
)

// parseSQLServerUUID converts SQL Server's mixed-endian UUID byte format to standard UUID
// SQL Server stores the first 3 groups (4 bytes, 2 bytes, 2 bytes) in little-endian
// while the last 2 groups (2 bytes, 6 bytes) remain in big-endian
func parseSQLServerUUID(b []byte) (uuid.UUID, error) {
	if len(b) != 16 {
		return uuid.UUID{}, fmt.Errorf("invalid UUID byte length: %d", len(b))
	}

	// Create a new byte array with corrected byte order
	corrected := make([]byte, 16)

	// Reverse first 4 bytes (Data1 - DWORD)
	corrected[0] = b[3]
	corrected[1] = b[2]
	corrected[2] = b[1]
	corrected[3] = b[0]

	// Reverse next 2 bytes (Data2 - WORD)
	corrected[4] = b[5]
	corrected[5] = b[4]

	// Reverse next 2 bytes (Data3 - WORD)
	corrected[6] = b[7]
	corrected[7] = b[6]

	// Copy remaining 8 bytes as-is (Data4 - byte array)
	copy(corrected[8:], b[8:])

	return uuid.FromBytes(corrected)
}

// toSQLServerUUID converts a standard UUID to SQL Server's mixed-endian byte format
// This is the reverse of parseSQLServerUUID - converts from standard to SQL Server format
func toSQLServerUUID(u uuid.UUID) []byte {
	b := u[:]
	sqlBytes := make([]byte, 16)

	// Reverse first 4 bytes (Data1 - DWORD)
	sqlBytes[0] = b[3]
	sqlBytes[1] = b[2]
	sqlBytes[2] = b[1]
	sqlBytes[3] = b[0]

	// Reverse next 2 bytes (Data2 - WORD)
	sqlBytes[4] = b[5]
	sqlBytes[5] = b[4]

	// Reverse next 2 bytes (Data3 - WORD)
	sqlBytes[6] = b[7]
	sqlBytes[7] = b[6]

	// Copy remaining 8 bytes as-is (Data4 - byte array)
	copy(sqlBytes[8:], b[8:])

	return sqlBytes
}
func TransformUUID(u uuid.UUID) (uuid.UUID, error) {
	// If it's already version 4, return as is
	if u.Version() == 4 {
		return u, nil
	}

	b := toSQLServerUUID(u)
	// Set version to 4 (bits 12–15)
	//b[6] = (b[6] & 0x0F) | (4 << 4)
	// Set variant to RFC 4122 (bits 6–7 = 10)
	//	b[8] = (b[8] & 0x3F) | 0x80

	return uuid.FromBytes(b)
}
func TransformUUIDToSQLServerV2(u uuid.UUID) (uuid.UUID, error) {
	// If it's already version 4, return as is
	b := toSQLServerUUID(u)
	// Set version to 4 (bits 12–15)
	//b[6] = (b[6] & 0x0F) | (4 << 4)
	// Set variant to RFC 4122 (bits 6–7 = 10)
	//	b[8] = (b[8] & 0x3F) | 0x80

	return uuid.FromBytes(b)
}
