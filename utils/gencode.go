package utils
import(
	"fmt"
	"math/rand"
	"time"
)
// Generate a 6-digit verification code
func GenerateVerificationToken() string {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator
	code := rand.Intn(900000) + 100000 // Generate a random number between 100000 and 999999
	return fmt.Sprintf("%06d", code)   // Format it as a 6-digit string
}