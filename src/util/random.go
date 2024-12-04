package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

var (
	once   sync.Once
	rng    *rand.Rand
	source rand.Source
)

// init initializes the random number generator only once
func init() {
	once.Do(func() {
		source = rand.NewSource(time.Now().UnixNano())
		rng = rand.New(source)
	})
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rng.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomName generates a random owner name
func RandomName() string {
	return RandomString(6)
}

// RandomEmail generates a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomUUID generates a random UUID
func RandomUUID() uuid.UUID {
	return uuid.New()
}

// RandomPhone generates a random phone number in E.164 format
func RandomPhone() string {
	// Country codes list - for simplicity, we'll use only a few common ones
	countryCodes := []int{1, 44, 61, 62, 86, 91} // USA, UK, Australia, Indonesia, China, India

	// Select a random country code
	countryCode := countryCodes[rng.Intn(len(countryCodes))]

	// Generate a random subscriber number.
	// Since E.164 allows up to 15 digits total, we'll use 10 for simplicity to fit within typical phone number lengths
	subscriberNumber := generateSubscriberNumber(10)

	return fmt.Sprintf("+%d%s", countryCode, subscriberNumber)
}

// generateSubscriberNumber generates a random string of digits
func generateSubscriberNumber(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteString(strconv.Itoa(rng.Intn(10)))
	}
	return sb.String()
}

// RandomBirthDate generates a random birth date in time.Time format
// It assumes a range from 100 years ago to 18 years ago (assuming legal adulthood at 18)
func RandomBirthDate() time.Time {
	// Current year
	currentYear := time.Now().Year()

	// Define the range for birth years
	// Assuming the youngest adult is 18, and oldest is 100
	minAge := 18
	maxAge := 100
	birthYear := currentYear - rng.Intn(maxAge-minAge+1) - minAge

	// Random month
	month := time.Month(rng.Intn(12) + 1)

	// Days in the month can vary, so we'll get the max days for the chosen month/year
	maxDays := time.Date(birthYear, month+1, 0, 0, 0, 0, 0, time.UTC).Day()

	// Random day within the month
	day := rng.Intn(maxDays) + 1

	// Create the date
	return time.Date(birthYear, month, day, 0, 0, 0, 0, time.UTC)
}

// RandomGender returns a random gender string, either "male" or "female"
func RandomGender() string {
	genders := []string{"male", "female"}
	return genders[rng.Intn(len(genders))]
}
