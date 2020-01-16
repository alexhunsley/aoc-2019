package main

import (
	"bufio"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	"strconv"
	"strings"
	//"github.com/BenLubar/memoize"
)


type ingredient struct {
	name string
	count int
}

type recipe struct {
	product ingredient
	ingredients []ingredient
}

func readInput(filename string) []string {
	file, err := os.Open(filename)

	if err != nil {
		panic("failed opening input file")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines = []string{}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func addBaseRequirement(baseRequirements *map[string]int, name string, count int) {
	currCount, ok := (*baseRequirements)[name];

	if !ok {
		(*baseRequirements)[name] = count
		return
	}
	(*baseRequirements)[name] = currCount + count
}

// no RHS is repeated in the input,
// so we don't worry about that
func main() {
	//lines := readInput("input2.txt")
	lines := readInput("input3.txt")
	//lines := readInput("input2.txt")
	//lines := readInput("input.txt")

	allRecipes := parseRecipes(lines)

	spew.Dump(allRecipes)

	baseRequirements := map[string]int{}

	findOreRequirement("FUEL", 1, allRecipes, &baseRequirements)
	//findOreRequirement("A", 1, allRecipes, &baseRequirements)

	fmt.Println("base reqs: ", baseRequirements)
	fmt.Println("=============================================")

	totalOreCost := findOreCost(allRecipes, &baseRequirements)
	fmt.Println("Day 14 part 1: total ore cost: ", totalOreCost)
}

func findOreCost(recipes map[string]recipe, baseRequirements *map[string]int) int {
	totalOreCost := 0

	for chemical, amountNeeded := range *baseRequirements {
		fmt.Println(chemical, amountNeeded)

		recipe := recipes[chemical]
		fmt.Println("recipe: ", recipe)

		fmt.Println("orig count: ", amountNeeded)
		//round up to nearest higher whole ingredient count
		//remainder := amountNeeded % recipe.product.count
		//if remainder > 0 {
		//	amountNeeded += recipe.product.count - (amountNeeded % recipe.product.count)
		//	//fmt.Println("new count: ", count)
		//}

		oldAN := amountNeeded
		amountNeeded := countAtNextHighestMultiple(amountNeeded, recipe.product.count)
		//cost := amountNeeded * recipe.ingredients[0].count
		cost := amountNeeded

		//multiplier := amountNeeded / recipe.product.count
		//cost := multiplier * recipe.ingredients[0].count
		fmt.Println("Found pre-amountNeeed, post-amountNeeded, factor: ", oldAN, amountNeeded, recipe.product.count)

		totalOreCost += cost
	}
	return totalOreCost
}

func findOreRequirement(itemName string, count int, recipes map[string]recipe, baseRequirements *map[string]int) {
	recipe := recipes[itemName]

	fmt.Println("----> looking at requirement of ", count, itemName, " (curr recipe: ", recipe)
	for _, ingredient := range recipe.ingredients {
		// we need count, and the item is made in multiples of ingredient.count.
		// adjust count to be a whole multiple of ingredient.count.

		if ingredient.name == "ORE" {
			//addBaseRequirement(baseRequirements, itemName, count * ingredient.count / recipe.product.count)
			//addBaseRequirement(baseRequirements, itemName, count * recipe.product.count)
			//addBaseRequirement(baseRequirements, itemName, count / recipe.product.count)
			addBaseRequirement(baseRequirements, itemName, count)
			return
		}

		count = countAtNextHighestMultiple(count, recipe.product.count)
 		findOreRequirement(ingredient.name, count * ingredient.count / recipe.product.count, recipes, baseRequirements)
	}
}

func countAtNextHighestMultiple(count int, factor int) int {
	fmt.Println("countAtNextHighestMultiple: count, factor = ", count, factor)
	remainder := count % factor
	if remainder == 0 {
		fmt.Println("countAtNextHighestMultiple:     .... no adjustment")
		return count
	}
	count += factor - remainder
	fmt.Println("countAtNextHighestMultiple:     .... adjustment to ", count)
	return count
}

func parseRecipes(lines []string) map[string]recipe {
	allRhs := []string{}

	// map from recipe name (e.g. ABC) to ingredient --> ingredients
	// e.g.:
	//      str(ABC) -> ingredient[ABC, 5] : [ingredient[DEF, 2], ingredient[XYZ, 10]]
	allRecipes := map[string]recipe{}

	for _, line := range lines {

		leftRight := strings.Split(line, " => ")

		leftSidePart := strings.Split(leftRight[0], ", ")

		rightSidePart := leftRight[1]
		rhsIngredient := stringToIngredient(rightSidePart)

		allRhs = append(allRhs, rhsIngredient.name)

		//fmt.Println("rhs :", rhsIngredient)
		//fmt.Println(leftSidePart)

		allIngredients := []ingredient{}

		//fmt.Println("Splitting ONE LHS:")
		for _, leftPart := range leftSidePart {
			ing := stringToIngredient(leftPart)
			allIngredients = append(allIngredients, ing)
			//fmt.Println("One ingredient: ", ing)
		}

		recipe := recipe{
			product:     rhsIngredient,
			ingredients: allIngredients,
		}

		// add to map from ingredient name to the components
		allRecipes[rhsIngredient.name] = recipe
	}
	return allRecipes
}

func stringToIngredient(leftPart string) ingredient {
	//fmt.Println("---- splitting this: ", leftPart)
	leftPartSplit := strings.Split(leftPart, " ")

	count, _ := strconv.Atoi(leftPartSplit[0])

	ing := ingredient{
		name:  leftPartSplit[1],
		count: count,
	}
	return ing
}
