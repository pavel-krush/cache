package main

import (
	"bufio"
	"fmt"
	"github.com/pavel-krush/cache/v2/lru/queue"
	"os"
)

func main() {
	q := queue.New(10)
	_=q

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("enter command: ")
		if !scanner.Scan() {
			break
		}

		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		switch line[0] {
		case 'p': // push
			if len(line) != 2 {
				fmt.Printf("bad push command\n")
				break
			}
			key := fmt.Sprintf("%c", line[1])
			q.Push(key)
			q.DebugPrint()
		case 'd': // delete
			if len(line) != 2 {
				fmt.Printf("bad delete command\n")
			}
			key := fmt.Sprintf("%c", line[1])
			q.Delete(key)
			q.DebugPrint()
		case 's': // shift
			key, found := q.Shift()
			if !found {
				fmt.Printf("<empty>\n")
			}

			fmt.Printf("key: \"%s\"\n", key)
			q.DebugPrint()
		//case 'k': // peek
		//	key, found := q.Peek()
		//	if !found {
		//		fmt.Printf("<empty>\n")
		//	}
		//
		//	fmt.Printf("key: \"%s\"\n", key)
		default:
		}
	}
}
