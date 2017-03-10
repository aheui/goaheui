package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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

const (
	SStack = iota
	SQueue
	SPipe
)

type Char struct {
	Lead  rune
	Vowel rune
	Tail  rune
}

func (c Char) isComplete() bool {
	empty := rune(' ')

	return c.Lead != empty && c.Vowel != empty
}

type Storage struct {
	StorageType int
	Memory      []int
}

func (s *Storage) pop() (int, bool) {
	if len(s.Memory) == 0 {
		return 0, false
	}

	var x int
	var xs []int

	if s.StorageType == SStack || s.StorageType == SPipe {
		x, xs = s.Memory[len(s.Memory)-1], s.Memory[:len(s.Memory)-1]
		s.Memory = xs
	} else {
		x, xs = s.Memory[0], s.Memory[1:]
		s.Memory = xs
	}

	return x, true
}

func (s *Storage) push(val int) {
	if s.StorageType == SStack || s.StorageType == SPipe {
		s.Memory = append(s.Memory, val)
	} else {
		s.Memory = append([]int{val}, s.Memory...)
	}
}

func (s Storage) peek() int {
	if s.StorageType == SStack || s.StorageType == SPipe {
		return s.Memory[len(s.Memory)-1]
	}

	return s.Memory[0]
}

type Machine struct {
	Codespace      [][]Char
	CurrentStorage *Storage
	xPos           int
	yPos           int
	dx             int
	dy             int
	terminated     bool
}

func (m *Machine) reverseCursorX() {
	m.dx = -m.dx
}

func (m *Machine) reverseCursorY() {
	m.dy = -m.dy
}

// reverseCursor reverses the direction of the momentum
func (m *Machine) reverseCursor() {
	m.reverseCursorX()
	m.reverseCursorY()
}

// moveCursor updates the position of the cursor based on the current 'momentum'
func (m *Machine) moveCursor() {
	m.xPos += m.dx
	m.yPos += m.dy

	xMaxPos := len(m.Codespace[0])
	yMaxPos := len(m.Codespace)

	if m.xPos > xMaxPos {
		m.xPos = m.xPos - xMaxPos
	} else if m.xPos < 0 {
		m.xPos = xMaxPos - m.dx
	}

	if m.yPos > yMaxPos {
		m.yPos = m.yPos - yMaxPos
	} else if m.yPos < 0 {
		m.yPos = yMaxPos - m.dy
	}
}

var stacks []Storage
var queue Storage
var pipe Storage
var KOREAN_OFFSET rune = 0xAC00
var DEBUG string

// Initialize stacks, queue, and pipe
func init() {
	// Set debug flag
	DEBUG, _ = os.LookupEnv("DEBUG")

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

	pipe = Storage{
		StorageType: SPipe,
		Memory:      []int{},
	}
}

// validateAheuiChar checks if a rune is a valid Aheui character
func validateAheuiChar(c rune) bool {
	return c >= 0xAC00 && c <= 0xD7A3
}

// makeChar deconstructs the rune corresponding to the Korean alphabet
// into lead, vowel, and tail sounds, and returns Char
func makeChar(c rune) Char {
	if !validateAheuiChar(c) {
		return Char{
			Lead:  rune(' '),
			Vowel: rune(' '),
			Tail:  rune(' '),
		}
	}

	codeNum := c - KOREAN_OFFSET

	tailNum := codeNum % 28
	vowelNum := (codeNum / 28) % 21
	leadNum := codeNum / 28 / 21

	lead := leadSounds[leadNum]
	vowel := vowelSounds[vowelNum]

	tail := rune(' ')
	if tailNum > 0 {
		tail = tailSounds[tailNum]
	}

	return Char{
		Lead:  lead,
		Vowel: vowel,
		Tail:  tail,
	}
}

func sanitizeCode(code string) string {
	reg := regexp.MustCompile("\n$")
	res := reg.ReplaceAllString(code, "")

	return strings.Trim(res, " ")
}

// initCodespace returns a CodeSpace based on the code string
func initCodespace(code string) ([][]Char, error) {
	lines := strings.Split(sanitizeCode(code), "\n")
	codeSpace := make([][]Char, len(lines))

	for lineIdx, line := range lines {
		codeSpace[lineIdx] = make([]Char, len(line)/3)
		for charIdx, char := range line {
			// WHY: why is index multiple of 3? (e.g. 3,6,9,...)
			codeSpace[lineIdx][charIdx/3] = makeChar(char)
		}
	}

	if DEBUG != "" {
		fmt.Printf("%+v", codeSpace)
	}

	return codeSpace, nil
}

// step evaluates the current Character and moves the cursor accordingly
func (m *Machine) step() int {
	currentChar := m.Codespace[m.yPos][m.xPos]

	if !currentChar.isComplete() {
		m.moveCursor()
		return 0
	}

	switch currentChar.Vowel {
	case 'ㅏ':
		m.dx = 1
		m.dy = 0
	case 'ㅓ':
		m.dx = -1
		m.dy = 0
	case 'ㅜ':
		m.dx = 0
		m.dy = 1
	case 'ㅗ':
		m.dx = 0
		m.dy = -1
	case 'ㅑ':
		m.dx = 2
		m.dy = 0
	case 'ㅕ':
		m.dx = -2
		m.dy = 0
	case 'ㅠ':
		m.dx = 0
		m.dy = 2
	case 'ㅛ':
		m.dx = 0
		m.dy = -2
	case 'ㅡ':
		if m.dy != 0 {
			m.reverseCursorY()
			break
		}
	case 'ㅣ':
		if m.dx != 0 {
			m.reverseCursorX()
		}
	case 'ㅢ':
		m.reverseCursor()
	default:
		//noop
	}

	switch currentChar.Lead {
	case 'ㅇ':
		// noop
		break
	case 'ㅎ':
		m.terminated = true

		popped, ok := m.CurrentStorage.pop()
		if ok {
			return popped
		}

		return 0
	case 'ㄷ':
		a, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}
		b, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		m.CurrentStorage.push(a + b)
		break
	case 'ㄸ':
		a, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}
		b, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		m.CurrentStorage.push(a * b)
		break
	case 'ㄴ':
		a, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}
		b, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		m.CurrentStorage.push(b / a)
		break
	case 'ㅌ':
		a, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}
		b, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		m.CurrentStorage.push(b - a)
		break
	case 'ㄹ':
		a, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}
		b, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		m.CurrentStorage.push(b % a)
		break
	case 'ㅁ':
		popped, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		switch currentChar.Tail {
		case 'ㅇ':
			fmt.Printf("%d", popped)
		case 'ㅎ':
			fmt.Printf("%s", string(popped))
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

	case 'ㅃ':
		i := m.CurrentStorage.peek()
		m.CurrentStorage.push(i)

	case 'ㅍ':
		a, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}
		b, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		m.CurrentStorage.push(a)
		m.CurrentStorage.push(b)

	case 'ㅅ':
		switch currentChar.Tail {
		case 'ㅇ':
			m.CurrentStorage = &queue
			break
		case 'ㅎ':
			m.CurrentStorage = &pipe
			break
		default:
			stackIdx := stackIndices[currentChar.Tail]
			m.CurrentStorage = &stacks[stackIdx]
		}

	case 'ㅆ':
		popped, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		switch currentChar.Tail {
		case 'ㅇ':
			queue.push(popped)
			break
		case 'ㅎ':
			pipe.push(popped)
			break
		default:
			stackIdx := stackIndices[currentChar.Tail]
			stacks[stackIdx].push(popped)
		}

	case 'ㅈ':
		a, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}
		b, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		var res int
		if b > a {
			res = 1
		} else {
			res = 0
		}

		m.CurrentStorage.push(res)

	case 'ㅊ':
		popped, ok := m.CurrentStorage.pop()
		if !ok {
			m.reverseCursor()
			break
		}

		if popped == 0 {
			m.reverseCursor()
		}
	}

	m.moveCursor()

	return 0
}

func (m *Machine) run(code string) (int, error) {
	codeSpace, err := initCodespace(code)
	if err != nil {
		return 0, err
	}

	m.Codespace = codeSpace

	var res int
	terminatedFlag := false

	for !terminatedFlag {
		res = m.step()
		terminatedFlag = m.terminated
	}

	return res, nil
}

// readFile reads the content of a file as string
func readFile(filepath string) (string, error) {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func main() {
	_ = makeChar('ㅂ')
	if len(os.Args) < 2 {
		fmt.Println("Please provide filepath")
		os.Exit(1)
	}

	filepath := os.Args[1]

	content, err := readFile(filepath)

	if err != nil {
		panic(err)
	}

	machine := Machine{
		CurrentStorage: &stacks[0],
		xPos:           0,
		yPos:           0,
		dx:             0,
		dy:             1,
		terminated:     false,
	}

	res, err := machine.run(content)
	if err != nil {
		panic(err)
	}

	os.Exit(res)
}
