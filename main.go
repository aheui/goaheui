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

var stackIndices = map[rune]int{
	' ': 0,
	'ㄱ': 1,
	'ㄴ': 2,
	'ㄷ': 3,
	'ㄹ': 4,
	'ㅁ': 5,
	'ㅂ': 6,
	'ㅅ': 7,
	'ㅈ': 8,
	'ㅊ': 9,
	'ㅋ': 10,
	'ㅌ': 11,
	'ㅍ': 12,
	'ㄲ': 13,
	'ㄳ': 14,
	'ㄵ': 15,
	'ㄶ': 16,
	'ㄺ': 17,
	'ㄻ': 18,
	'ㄼ': 19,
	'ㄽ': 20,
	'ㄾ': 21,
	'ㄿ': 22,
	'ㅀ': 23,
	'ㅄ': 24,
	'ㅆ': 25,
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

func (s *Storage) pop() int {
	var x int
	var xs []int

	if s.StorageType == SStack {
		x, xs = s.Memory[len(s.Memory)-1], s.Memory[:len(s.Memory)-1]
		s.Memory = xs
	} else {
		x, xs = s.Memory[0], s.Memory[1:]
		s.Memory = xs
	}

	return x
}

func (s *Storage) push(val int) {
	if s.StorageType == SStack {
		s.Memory = append(s.Memory, val)
	} else {
		s.Memory = append([]int{val}, s.Memory...)
	}
}

type Machine struct {
	CurrentStorage Storage
	xPos           int
	yPos           int
	dx             int
	dy             int
	terminated     bool
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
		terminated:     false,
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

func (m Machine) step(codeSpace [][]Char) int {
	currentChar := codeSpace[m.yPos][m.xPos]

	switch currentChar.Lead {
	case 'ㅇ':
		// noop
		break
	case 'ㅎ':
		m.terminated = true

		if len(m.CurrentStorage.Memory) > 0 {
			return m.CurrentStorage.pop()
		}

		return 0
	case 'ㄷ':
		a, b := m.CurrentStorage.pop(), m.CurrentStorage.pop()
		m.CurrentStorage.push(a + b)
		break
	case 'ㄸ':
		a, b := m.CurrentStorage.pop(), m.CurrentStorage.pop()
		m.CurrentStorage.push(a * b)
		break
	case 'ㄴ':
		a, b := m.CurrentStorage.pop(), m.CurrentStorage.pop()
		m.CurrentStorage.push(a / b)
		break
	case 'ㅌ':
		a, b := m.CurrentStorage.pop(), m.CurrentStorage.pop()
		m.CurrentStorage.push(a - b)
		break
	case 'ㄹ':
		a, b := m.CurrentStorage.pop(), m.CurrentStorage.pop()
		m.CurrentStorage.push(a % b)
		break
	case 'ㅁ':
		popped := m.CurrentStorage.pop()

		switch currentChar.Tail {
		case 'ㅇ':
			fmt.Println(popped)
		case 'ㅎ':
			fmt.Println(string(popped))
		}

	case 'ㅂ':
		switch currentChar.Tail {
		case 'ㅇ':
			var i int
			fmt.Scanf("%d", &i)
			m.CurrentStorage.push(i)
			break
		case 'ㅎ':
			var i rune
			fmt.Scanf("%c", &i)
			m.CurrentStorage.push(int(i))
			break
		default:
			m.CurrentStorage.push(strokeCount[currentChar.Tail])
		}

	case 'ㅅ':
		switch currentChar.Tail {
		case 'ㅇ':
			m.CurrentStorage = queue
			break
		case 'ㅎ':
			// TODO: pipe
			break
		default:
			stackIdx := stackIndices[currentChar.Tail]
			m.CurrentStorage = stacks[stackIdx]
		}

	case 'ㅆ':
		popped := m.CurrentStorage.pop()

		switch currentChar.Tail {
		case 'ㅇ':
			queue.push(popped)
			break
		case 'ㅎ':
			// TODO: pipe
			break
		default:
			stackIdx := stackIndices[currentChar.Tail]
			stacks[stackIdx].push(popped)
		}

	case 'ㅈ':
		a, b := m.CurrentStorage.pop(), m.CurrentStorage.pop()

		var res int
		if b > a {
			res = 1
		} else {
			res = 0
		}

		m.CurrentStorage.push(res)

		//TODO
		//	case 'ㅊ':
		//		popped := m.CurrentStorage.pop()
		//
	}

	switch currentChar.Vowel {
	case 'ㅏ':
		m.xPos += 1
		m.dx = 1
		m.dy = 0
		break
	case 'ㅓ':
		m.xPos -= 1
		m.dx = -1
		m.dy = 0
		break
	case 'ㅜ':
		m.yPos += 1
		m.dx = 0
		m.dy = 1
		break
	case 'ㅗ':
		m.yPos -= 1
		m.dx = 0
		m.dy = -1
		break
	case 'ㅑ':
		m.xPos += 2
		m.dx = 2
		m.dy = 0
		break
	case 'ㅕ':
		m.xPos -= 2
		m.dx = -2
		m.dy = 0
		break
	case 'ㅠ':
		m.yPos += 2
		m.dx = 0
		m.dy = 2
		break
	case 'ㅛ':
		m.yPos -= 2
		m.dx = 0
		m.dy = -2
		break
	case 'ㅡ':
		if m.dy != 0 {
			m.yPos = m.yPos - m.dy
			m.dy = -m.dy
		}
	case 'ㅣ':
		if m.dx != 0 {
			m.xPos = m.xPos - m.dx
			m.dx = -m.dx
		}
	case 'ㅢ':
		m.xPos = m.xPos - m.dx
		m.yPos = m.yPos - m.dy
		m.dy = -m.dy
		m.dx = -m.dx
	default:
		//noop
	}

	return 0
}

func (m Machine) run(codeSpace [][]Char) int {
	var res int
	if !m.terminated {
		res = m.step(codeSpace)
	}

	return res
}

func main() {
	var codeSpace = initCodespace(helloWorld)

	machine.run(codeSpace)
}
