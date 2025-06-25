package utils

import "strings"

// GenerateAbbreviation creates an acronym/abbreviation from the given full name.
// Example: "Abubakar Tafawa Balewa University" -> "ATBU"
func GenerateAbbreviation(name string) string {
    if name == "" {
        return ""
    }
    words := strings.Fields(name)
    abbr := make([]rune, 0, len(words))
    for _, w := range words {
        // Skip empty strings or words that are just punctuation
        for _, r := range w {
            if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
                abbr = append(abbr, r)
                break
            }
        }
    }
    return strings.ToUpper(string(abbr))
} 