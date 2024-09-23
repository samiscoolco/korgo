package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// town constants
var prefixes = []string{"Cyr", "Kor", "Mal", "Sor", "Tal", "Lun", "Tar", "Elin", "Gon", "Fen", "Zar", "Vel", "Tor", "Dra", "Nim", "Sed", "Har"}
var middles = []string{"mar", "lor", "thil", "nex", "dor", "wen", "ven", "kor", "zil", "var", "mel", "ril", "san", "jel", "tan", "quor"}
var suffixes = []string{"ia", "ton", "dale", "mere", "burg", "field", "holm", "wood", "stead", "view", "crest", "haven", "brook", "ridge", "port", "cliff"}

func GenerateTownName(seed int) string {
	// Seed the random number generator
	r := rand.New(rand.NewSource(int64(seed)))

	// Randomly select parts
	prefix := prefixes[r.Intn(len(prefixes))]
	middle := middles[r.Intn(len(middles))]
	suffix := suffixes[r.Intn(len(suffixes))]

	// Combine parts to create the town name
	townName := fmt.Sprintf("%s%s%s", prefix, middle, suffix)
	return townName
}

// return random int between 1-6
func d6() int {
	return rand.Intn(6) + 1
}

// clear screen
func cls() {
	fmt.Print("\033[H\033[2J")
}

// Structure of an encounter
type encounter struct {
	nm string
	ds string
	th int
	dm int
	gp int
}

// Structure of an item
type item struct {
	name        string
	desc        string
	stat_num    int
	stat_change int
	stat_mod    int
	buyprice    int
}

func saveFile(name string, statdata []int) {
	var sb strings.Builder
	sb.WriteString(name)
	sb.WriteString(" ")
	for _, v := range statdata {
		sb.WriteString(strconv.Itoa(v))
		sb.WriteString(" ")
	}

	var strdata string = sb.String()
	data := []byte(strdata)

	err := os.WriteFile("save.kgo", data, 0644)
	if err != nil {
		fmt.Println("save unsuccessful.")
	}
	fmt.Println("save successful!")

}

func main() {
	var input string

	//player stats the weird and the wild
	var running bool = true
	var in_encounter bool = false
	var hero string

	// 0:hp, 1:gold, 2:trap_dodge, 3:flees, 4:shield, 5: atk, 6:armor, 7:heals, 8: kills
	// update this whenevear ading a stat pls

	var stat_names = []string{"HP", "GOLD", "ROPES", "TRAPS", "SHIELD", "ATK", "ARMOR", "HEALS", "KILLS", "PLAYER_LEVEL", "PLAYER_MAXLEVEL"}
	var stat_maxes = []int{20, -1, 3, 3, 1, 1, 1, 4, -1, -1, -1}
	var stats = []int{20, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1}
	var hp *int = &stats[0]
	var gold *int = &stats[1]
	var escape *int = &stats[2]
	var flees *int = &stats[3]
	var shields *int = &stats[4]
	var atk *int = &stats[5]
	var def *int = &stats[6]
	var heals *int = &stats[7]
	var kills *int = &stats[8]
	var player_level *int = &stats[9]
	//var player_level_max *int = &stats[10]

	var stat_total int = len(stats)

	var roll int = 0

	var action string
	var shopmode bool

	var current_encounter *encounter

	var encounter_list = []encounter{
		encounter{
			nm: "Pit Trap",
			ds: "a pit full of unforgiving wooden spikes..",
			th: -1,
			dm: 3,
			gp: 0,
		},
		encounter{
			nm: "Skeleton",
			ds: "the reanimated hollow interior of a previous adventurer.",
			th: 2,
			dm: 2,
			gp: 1,
		},
		encounter{
			nm: "Goblin",
			ds: "a fearsome green abomination.",
			th: 3,
			dm: 2,
			gp: 2,
		},
		encounter{
			nm: "Wraith",
			ds: "an ethereal whisper of malice.",
			th: 5,
			dm: 2,
			gp: 4,
		},
		encounter{
			nm: "Ogre",
			ds: "a towering mountain of meat and hatred with eyes.",
			th: 4,
			dm: 4,
			gp: 7,
		},
		encounter{
			nm: "Demon",
			ds: "a pestilent denison of another plane of existence.",
			th: 4,
			dm: 6,
			gp: 10,
		},
	}

	var shop_items = []item{
		item{
			name:        "Rope",
			desc:        "Avoid trap once",
			stat_num:    2,
			stat_change: 1,
			stat_mod:    0,
			buyprice:    2,
		},
		item{
			name:        "Caltrops",
			desc:        "Flee foe once",
			stat_num:    3,
			stat_change: 1,
			stat_mod:    0,
			buyprice:    2,
		},
		item{
			name:        "Shield",
			desc:        "Avoid Lethal Damage Once",
			stat_num:    4,
			stat_change: 1,
			stat_mod:    0,
			buyprice:    3,
		},
		item{
			name:        "Potion",
			desc:        "Heal D6 HP",
			stat_num:    7,
			stat_change: 1,
			stat_mod:    0,
			buyprice:    4,
		},
		item{
			name:        "Sword",
			desc:        "+1 to Defeat Foe",
			stat_num:    5,
			stat_change: 1,
			stat_mod:    0,
			buyprice:    5,
		},
		item{
			name:        "Armor",
			desc:        "-1 DMG to self",
			stat_num:    6,
			stat_change: 1,
			stat_mod:    0,
			buyprice:    10,
		},
	}

	shopmode = true

	//main menu
mm:
	cls()
	fmt.Print("Welcome to KorGO\n(n)ew game, (l)oad game, (q)uit\n")
	fmt.Print("Choice: \n")
	fmt.Scanln(&input)
	switch input {
	case "l":
		// load game
		dat, _ := os.ReadFile("save.kgo")
		fmt.Print(string(dat))
		svdata := strings.Split(string(dat), " ")
		for i, v := range svdata {
			if i == 0 {
				hero = v
				continue
			} else {
				if i-1 < stat_total {
					stats[i-1], _ = strconv.Atoi(v)
				}

			}
		}
	case "q":
		//quit game
		return
	case "n":
		fmt.Print("Who goes there?\n")
		fmt.Print("Name: \n")
		fmt.Scanln(&hero)
	default:
		//call me pathetic, call me what you will
		goto mm
	}

	for running {
		if shopmode {
			//clear screen
			cls()
			//shop screen
			fmt.Print(hero, " LVL ", *player_level, " HP: ", *hp, "  GOLD: ", *gold, "\n")
			fmt.Print("You are shopping in the town of ", GenerateTownName(*player_level), "\nItems: \n")
			for i, v := range shop_items {
				fmt.Print(i, " - ", v.name, " ", v.desc, " ", v.buyprice, "GP\n")
			}

			//autosave
			saveFile(hero, stats)

			//input
			fmt.Print("(0-9) buy item, (e)xplore, (i)nventory, (l)evel up, (t)ravel, (q)uit\n")
			fmt.Print("Choice: \n")
			fmt.Scanln(&input)

			i, err := strconv.Atoi(input)
			if err != nil {
				if input == "q" {
					saveFile(hero, stats)
					return
				}
				if input == "i" {
					cls()
					//print out inventory
					fmt.Print("Inventory/Stats for ", hero, "\n")
					for i := range stats {
						fmt.Print(stat_names[i], " - ", stats[i], "\n")
					}
					fmt.Print("Press Enter to Continue...\n")
					fmt.Scanln()
				}
				if input == "e" {
					shopmode = false
				} else {
					fmt.Print("Unknown Command\n")
				}
			} else {
				if *gold >= shop_items[i].buyprice {
					if stats[shop_items[i].stat_num] < stat_maxes[shop_items[i].stat_num] {
						fmt.Print()
						fmt.Print("Buy ", shop_items[i].name, " for ", shop_items[i].buyprice, "GP?", "\n")
						fmt.Print("(y)es, (n)o\n")
						fmt.Scanln(&input)
						if input == "y" {
							*gold = *gold - shop_items[i].buyprice
							stats[shop_items[i].stat_num] += shop_items[i].stat_change
						}
					} else {
						fmt.Print("You can hold no more ", stat_names[shop_items[i].stat_num], "!\n")
						fmt.Print("Press Enter to Continue...\n")
						fmt.Scanln()
					}

				} else {
					fmt.Print("You do not have enough GP!\n")
					fmt.Print("Press Enter to Continue...\n")
					fmt.Scanln()
				}

			}
		} else {
			for input != "t" && running {
				cls()
				fmt.Print(hero, " HP: ", *hp, "  GOLD: ", *gold, "\n")
				fmt.Print("You are exploring PLACENAME\n(c)ontinue, (t)own\nAction? \n")
				fmt.Scanln(&input)
				if input == "t" || input != "c" {
					if input == "t" {
						shopmode = true
					}
					continue
				}

				cls()

				roll = d6()
				current_encounter = &encounter_list[roll-1]
				fmt.Print("Before you is ", current_encounter.ds, "\n")
				fmt.Print("You encountered a(n) ", current_encounter.nm, "(", roll, ")... \nPress Enter to Continue...\n")
				fmt.Scanln()
				in_encounter = true

				//trap
				if current_encounter.th == -1 {
					in_encounter = false
					if *escape > 0 {
						fmt.Print("Use a tool to avoid the ", current_encounter.nm, "? \n(y)es,(n)o\n")
						fmt.Print("Choice: \n")
						fmt.Scanln(&input)
						if input == "y" {
							*escape = *escape - 1
							fmt.Print("You successfully avoid the ", current_encounter.nm, "!\n")
							fmt.Print("Press Enter to Continue...\n")
							fmt.Scanln()
							continue
						}
					}

					fmt.Print("The ", current_encounter.nm, " rends you for ", current_encounter.dm, " damage!\n")
					fmt.Print("Press Enter to Continue...\n")

					*hp = *hp - current_encounter.dm
					if *hp <= 0 {
						running = false
						break
					}
					fmt.Scanln()
					continue
				}

				for in_encounter && running {
					cls()
					fmt.Print(hero, " HP: ", *hp, "  GOLD: ", *gold, "\n")
					fmt.Print(current_encounter.nm, "(", current_encounter.th, current_encounter.dm, current_encounter.gp, ")\n")
					fmt.Print("(a)ttack, (f)lee, (h)eal\n")
					fmt.Print("Action? \n")
					fmt.Scanln(&action)
					switch action {
					case "a":
						roll = d6() + *atk
						fmt.Print("Attack roll ", roll, "\n")
						if roll > current_encounter.th {
							fmt.Print("Enemy slain, gain ", current_encounter.gp, " gold.\n")
							*gold = *gold + current_encounter.gp
							*kills = *kills + 1
							in_encounter = false
						} else {
							fmt.Print("You missed, you are hurt for ", (current_encounter.dm - *def), "\n")
							*hp = *hp - (current_encounter.dm - *def)
							if *hp <= 0 {
								if *shields > 0 {
									*shields = *shields - 1
									fmt.Print("Your shield breaks and saves your life.\n")
									*hp = d6()
								} else {
									in_encounter = false
									running = false
									break
								}
							}

						}
					case "f":
						if *flees > 0 {
							fmt.Print("Flee success\n")
							*flees = *flees - 1
							in_encounter = false
						} else {
							fmt.Print("You need a trap.\n")
						}
					case "h":
						if *heals > 0 {
							roll = d6()
							fmt.Print("Healed for ", roll, "\n")
							//dont forget max
							*hp = *hp + roll
							*heals = *heals - 1
						} else {
							fmt.Print("You have no heals.\n")
						}
					default:
						fmt.Print("Invalid Action\n")
					}
					fmt.Print("Press Enter to Continue...\n")
					fmt.Scanln()
				}
			}
		}
	}
	cls()
	fmt.Print("Your journey has come to an end.\nStats: \n")
	//print out inventory
	fmt.Print(hero, "\n")
	for i := range stats {
		fmt.Print(stat_names[i], " - ", stats[i], "\n")
	}
	fmt.Scanln()
	fmt.Print("Goodbye!")
	return
}
