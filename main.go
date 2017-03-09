package main

import (
	"fmt"
	"strings"
)

var strokeCount = map[rune]int{
	'ㄱ': 2,
	'ㄴ': 2,
	'ㄷ': 3,
	'ㄹ': 5,
	'ㅁ': 4,
	'ㅂ': 4,
	'ㅅ': 2,
	'ㅈ': 3,
	'ㅊ': 4,
	'ㅋ': 3,
	'ㅌ': 4,
	'ㅍ': 4,
	'ㄲ': 4,
	'ㄳ': 4,
	'ㄵ': 5,
	'ㄶ': 5,
	'ㄺ': 7,
	'ㄻ': 9,
	'ㄼ': 9,
	'ㄽ': 7,
	'ㄾ': 9,
	'ㄿ': 9,
	'ㅀ': 8,
	'ㅄ': 6,
	'ㅆ': 4,
}

var leadSounds = []rune{
	'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ',
	'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ',
	'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ',
	'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ',
}

var vowelSounds = []rune{
	'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ', 'ㅓ',
	'ㅔ', 'ㅕ', 'ㅖ', 'ㅗ', 'ㅘ',
	'ㅙ', 'ㅚ', 'ㅛ', 'ㅜ', 'ㅝ',
	'ㅞ', 'ㅟ', 'ㅠ', 'ㅡ', 'ㅢ',
	'ㅣ',
}

var tailSounds = []rune{
	' ', 'ㄱ', 'ㄲ', 'ㄳ', 'ㄴ',
	'ㄵ', 'ㄶ', 'ㄷ', 'ㄹ', 'ㄺ',
	'ㄻ', 'ㄼ', 'ㄽ', 'ㄾ', 'ㄿ',
	'ㅀ', 'ㅁ', 'ㅂ', 'ㅄ', 'ㅅ',
	'ㅆ', 'ㅇ', 'ㅈ', 'ㅊ', 'ㅋ',
	'ㅌ', 'ㅍ', 'ㅎ',
}

var helloWorld = `밤밣따빠밣밟따뿌
빠맣파빨받밤뚜뭏
돋밬탕빠맣붏두붇
볻뫃박발뚷투뭏붖
뫃도뫃희멓뭏뭏붘
뫃봌토범더벌뿌뚜
뽑뽀멓멓더벓뻐뚠
뽀덩벐멓뻐덕더벅`

const (
	SStack = iota
	SQueue
)

type Char struct {
	Lead  rune
	Vowel rune
	Tail  rune
}

var storages []Storage
var KOREAN_OFFSET rune = 0xAC00

type Storage struct {
	StorageType int
	Memory      []int
}

func (s Storage) pop() int {
	if s.StorageType == SStack {
		x, xs := s.Memory[len(s.Memory)-1], s.Memory[:len(s.Memory)-1]
		s.Memory = xs
		return x
	}

	x, xs := s.Memory[0], s.Memory[1:]
	s.Memory = xs
	return x
}

func (s Storage) push(val int) {
	if s.StorageType == SStack {
		s.Memory = append(s.Memory, val)
	}

	s.Memory = append([]int{val}, s.Memory...)
}

type Machine struct {
	CurrentStorage Storage
	xPos           int
	yPos           int
	dx             int
	dy             int
}

var stacks []Storage
var queue Storage
var machine Machine

func init() {
	for i := 0; i < 26; i++ {
		stack := Storage{
			StorageType: SStack,
			Memory:      []int{},
		}

		stacks = append(stacks, stack)
	}

	queue = Storage{
		StorageType: SQueue,
		Memory:      []int{},
	}

	machine = Machine{
		CurrentStorage: stacks[0],
		xPos:           0,
		yPos:           0,
		dx:             0,
		dy:             1,
	}
}

func validateAheuiChar(c rune) bool {
	return c >= 0xAC00 && c <= 0xD7A3
}

func makeChar(c rune) Char {
	codeNum := c - KOREAN_OFFSET

	tailNum := codeNum % 28
	vowelNum := (codeNum / 28) % 21
	leadNum := codeNum / 28 / 21

	lead := leadSounds[leadNum]
	vowel := vowelSounds[vowelNum]
	var tail rune
	if tailNum > 0 {
		tail = tailSounds[tailNum]
	} else {
		tail = 0
	}

	return Char{
		Lead:  lead,
		Vowel: vowel,
		Tail:  tail,
	}
}

func initCodespace(input string) [][]Char {
	lines := strings.Split(input, "\n")

	codeSpace := make([][]Char, len(lines))

	for lineIdx, line := range lines {
		codeSpace[lineIdx] = make([]Char, len(lines))
		for charIdx, char := range line {
			// TODO: why is index multiple of 3? (e.g. 3,6,9,...)
			codeSpace[lineIdx][charIdx/3] = makeChar(char)
		}
	}

	return codeSpace
}

func (m Machine) run(codeSpace [][]Char) {
	currentChar := codeSpace[m.yPos][m.xPos]
	fmt.Println(currentChar)

	switch currentChar.Lead {
	case 'ㅇ':
		// noop
		break
	case 'ㅎ':
		// TODO: pop
		fmt.Println("Terminate")
	case 'ㄷ':

	case 'ㅁ':
		popped := m.CurrentStorage.pop()
		fmt.Println(popped)
		switch currentChar.Tail {
		case 'ㅇ':
			fmt.Println(popped)
		case 'ㅎ':
			fmt.Println(string(popped))
		}

	case 'ㅂ':
		m.CurrentStorage.push(1)
	}
}

func main() {
	var codeSpace = initCodespace("몽")

	machine.run(codeSpace)
}
