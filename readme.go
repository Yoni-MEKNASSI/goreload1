package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// 1. Vérifier si on a bien donné le fichier d'entrée et de sortie en argument
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run . <fichier_entree> <fichier_sortie>")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// 2. Lire le fichier d'entrée
	content, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du fichier :", err)
		return
	}

	// Convertir le contenu du fichier en chaîne de caractères (string)
	text := string(content)

	// 3. Appliquer les transformations sur le texte
	processedText := transformText(text)

	// 4. Écrire le résultat dans le fichier de sortie
	err = os.WriteFile(outputFile, []byte(processedText), 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier :", err)
		return
	}

	fmt.Println("Transformation terminée avec succès !")
}

// Fonction principale pour transformer le texte mot par mot
func transformText(text string) string {
	// Découper le texte en une liste de mots
	words := strings.Fields(text)
	var result []string

	for i := 0; i < len(words); i++ {
		word := words[i]

		// Gestion de (hex) : transforme le mot précédent de hexadécimal en décimal[cite: 1]
		if word == "(hex)" && len(result) > 0 {
			val, err := strconv.ParseInt(result[len(result)-1], 16, 64)
			if err == nil {
				result[len(result)-1] = strconv.FormatInt(val, 10)
			}
			continue
		}

		// Gestion de (bin) : transforme le mot précédent de binaire en décimal[cite: 1]
		if word == "(bin)" && len(result) > 0 {
			val, err := strconv.ParseInt(result[len(result)-1], 2, 64)
			if err == nil {
				result[len(result)-1] = strconv.FormatInt(val, 10)
			}
			continue
		}

		// Gestion de (up) : met le mot précédent en majuscules[cite: 1]
		if word == "(up)" && len(result) > 0 {
			result[len(result)-1] = strings.ToUpper(result[len(result)-1])
			continue
		}

		// Gestion de (low) : met le mot précédent en minuscules[cite: 1]
		if word == "(low)" && len(result) > 0 {
			result[len(result)-1] = strings.ToLower(result[len(result)-1])
			continue
		}

		// Gestion de (cap) : met une majuscule initiale au mot précédent[cite: 1]
		if word == "(cap)" && len(result) > 0 {
			result[len(result)-1] = capitalize(result[len(result)-1])
			continue
		}

		// Si ce n'est pas un modificateur, on ajoute simplement le mot
		result = append(result, word)
	}

	// Reconstruire le texte avec les mots transformés
	finalText := strings.Join(result, " ")

	// Appliquer les corrections de ponctuation et des articles a/an[cite: 2]
	finalText = fixPunctuation(finalText)
	finalText = fixArticles(finalText)

	return finalText
}

// Fonction utilitaire pour mettre une majuscule au début d'un mot
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
}

// Fonction pour corriger la ponctuation (collée avant, espacée après)[cite: 2]
func fixPunctuation(text string) string {
	punctuations := []string{".", ",", "!", "?", ":", ";"}
	for _, p := range punctuations {
		// Enlever l'espace avant la ponctuation
		text = strings.ReplaceAll(text, " "+p, p)
		// Ajouter un espace après la ponctuation
		text = strings.ReplaceAll(text, p, p+" ")
	}

	// Nettoyer les doubles espaces créés par inadvertance
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return strings.TrimSpace(text)
}

// Fonction pour transformer "a" en "an" devant une voyelle ou un 'h'[cite: 2]
func fixArticles(text string) string {
	words := strings.Fields(text)
	for i := 0; i < len(words)-1; i++ {
		w := strings.ToLower(words[i])
		if w == "a" {
			next := strings.ToLower(words[i+1])
			if len(next) > 0 {
				first := next[0]
				// Si le mot suivant commence par une voyelle ou un 'h'
				if first == 'a' || first == 'e' || first == 'i' || first == 'o' || first == 'u' || first == 'h' {
					if words[i] == "A" {
						words[i] = "An"
					} else {
						words[i] = "an"
					}
				}
			}
		}
	}
	return strings.Join(words, " ")
}
