package main

import "flag"

var GoReserved = []string{
	"break",        "default",      "func",         "interface",    "select",
	"case",         "defer",        "go",           "map",          "struct",
	"chan",         "else",         "goto",         "package",      "switch",
	"const",        "fallthrough",  "if",           "range",        "type",
	"continue",     "for",          "import",       "return",       "var",
	"bool",			"byte", 		"len", 			"open", 		"file", 
	"close", 		"load", 		"copy",
}

//This is the Java compiler for uct.
var Go bool

func init() {
	flag.BoolVar(&Go, "go", false, "Target Go")

	RegisterAssembler(GoAssembly, &Go, "go", "//")

	for _, word := range GoReserved {
		GoAssembly[word] = Reserved()
	}
}

var GoAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   "package main",
		Args: 1,
	},

	"FOOTER": Instruction{},

	"FILE": Instruction{
		Data: GoFile,
		Path: "/Stack.go",
	},

	"NUMBER": is("NewNumber(%s)", 1),
	"BIG": 	is("BigInt(`%s`)", 1),
	"SIZE":   is("%s.Len()", 1),
	"STRING": is("NewStringArray(%s)", 1),
	"ERRORS":  is("stack.ERROR", 1),

	"SOFTWARE": Instruction{
		Data:   "func main() { stack := &Stack{}; stack.Init();",
		Indent: 1,
	},
	"EXIT": Instruction{
		Indented:    1,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "os.Exit(stack.ERROR.ToInt())",
		},
	},

	"FUNCTION": is("func %s(stack *Stack) {", 1, 1),
	"RETURN": Instruction{
		Indented:    1,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return",
		},
	},
	
	"SCOPE": is(`stack.Relay(Pipe{Function:%s})`, 1),
	
	"EXE": is("%s.Exe(stack)", 1),

	"PUSH": is("stack.Push(%s)", 1),
	"PULL": is("%s := stack.Pull(); %s.Init()", 1),

	"PUT":   is("stack.Put(%s)", 1),
	"POP":   is("%s := stack.Pop()", 1),
	"PLACE": is("stack.ActiveArray = &%s", 1),
	"ARRAY":  is("var %s Array; stack.ActiveArray = &%s", 1),
	"RENAME": is("%s = *stack.ActiveArray", 1),

	"SHARE": is("stack.Share(%s)", 1),
	"GRAB":  is("%s := stack.Grab(); %s.Init()", 1),

	"RELAY": is("stack.Relay(%s)", 1),
	"TAKE":  is("%s := stack.Take(); %s.Init()", 1),

	"GET": is("%s := stack.Get()", 1),
	"SET": is("stack.Set(%s)", 1),

	"VAR": is("var %s Number", 1),

	"OPEN":   is("stack.Open()"),
	"LOAD":   is("stack.Load()"),
	"OUT":    is("stack.Out()"),
	"STAT":   is("stack.Info()"),
	"IN":     is("stack.In()"),
	"STDOUT": is("stack.Stdout()"),
	"STDIN":  is("stack.Stdin()"),
	"HEAP":  is("stack.Heap()"),

	"CLOSE": is("%s.Close()", 1),

	"LOOP":   is("for {", 0, 1),
	"BREAK":  is("break"),
	"REPEAT": is("}", 0, -1, -1),

	"IF":   is("if %s.True()  {", 1, 1),
	"ELSE": is("} else {", 0, 0, -1),
	"END":  is("}", 0, -1, -1),

	"RUN":  is("%s(stack)", 1),
	"DATA": is("var %s Array = %s;", 2),

	"FORK": is("go %s(stack.Copy())\n", 1),

	"ADD": is("%s.Add(%s, %s)", 3),
	"SUB": is("%s.Sub(%s, %s)", 3),
	"MUL": is("%s.Mul(%s, %s)", 3),
	"DIV": is("%s.Div(%s, %s)", 3),
	"MOD": is("%s.Mod(%s, %s)", 3),
	"POW": is("%s.Pow(%s, %s)", 3),

	"SLT": is("%s = %s.Slt(%s)", 3),
	"SEQ": is("%s = %s.Seq(%s)", 3),
	"SGE": is("%s = %s.Sge(%s)", 3),
	"SGT": is("%s = %s.Sgt(%s)", 3),
	"SNE": is("%s = %s.Sne(%s)", 3),
	"SLE": is("%s = %s.Sle(%s)", 3),

	"JOIN": is("%s = %s.Join(%s);", 3),
	"ERROR": is("stack.ERROR = %s;", 1),
}

//Edit this in a Java IDE.
const GoFile = `
//Compiled to Go with UCT (Universal Code Translator)
package main

import "math/big"
import "os"
import "io"
import "crypto/rand"
import "net"
import "strings"
import "strconv"
import "bufio"

var Networks_In = make(map[string]net.Listener)

//This is the Go stack implementation.
// It holds arrays for the 3 types:
//		Numbers
//		Arrays
//		Pipes
//
// It also holds the ERROR variable for the current thread.
// The currently active array is stored as ActiveArray.
type Stack struct {
	Numbers Array
	Arrays 	[]Array
	Pipes	[]Pipe
	
	ERROR Number
	ActiveArray *Array
	TheHeap []Array;
	HeapRoom []int
}

func (stack *Stack) Copy() (n *Stack) {
	n = new(Stack)
	n.Numbers = Array{}
	n.Numbers.Small = make([]byte, len(stack.Numbers.Small))
	copy(n.Numbers.Small, stack.Numbers.Small)
	n.Numbers.Big = make([]Number, len(stack.Numbers.Big))
	copy(n.Numbers.Big, stack.Numbers.Big)
	
	n.Arrays = make([]Array, len(stack.Arrays))
	for i := range stack.Arrays {
		n.Arrays[i] = Array{}
		n.Arrays[i].Small = make([]byte, len(stack.Arrays[i].Small))
		copy(n.Arrays[i].Small, stack.Arrays[i].Small)
		
		n.Arrays[i].Big = make([]Number, len(stack.Arrays[i].Big))
		copy(n.Arrays[i].Big, stack.Arrays[i].Big)
	}
	
	n.Pipes = make([]Pipe, len(stack.Pipes))
	copy(n.Pipes, stack.Pipes)
	
	return
}

func (stack *Stack) Array() Array {
	var array Array
	stack.ActiveArray = &array
	return array
}

func (z *Stack) Init() {

}

func (stack *Stack) Share(array Array) {
	stack.Arrays = append(stack.Arrays, array)
}
func (stack *Stack) Grab() (array Array) {
	array = stack.Arrays[len(stack.Arrays)-1]
	stack.Arrays = stack.Arrays[:len(stack.Arrays)-1]
	return
}

func (stack *Stack) Relay(pipe Pipe) {
	stack.Pipes = append(stack.Pipes, pipe)
}
func (stack *Stack) Take() (pipe Pipe) {
	pipe = stack.Pipes[len(stack.Pipes)-1]
	stack.Pipes = stack.Pipes[:len(stack.Pipes)-1]
	return
}

func (stack *Stack) Push(number Number) {
	if number.Int == nil && number.Small < 256 && number.Small >= 0 && stack.Numbers.Big == nil {
		stack.Numbers.Small = append(stack.Numbers.Small, byte(number.Small))
	} else {
		stack.Numbers.Grow()
		stack.Numbers.Big = append(stack.Numbers.Big, number)
	}
}
func (stack *Stack) Pull() (number Number) {
	if stack.Numbers.Big == nil {
		number.Small = int64(stack.Numbers.Small[len(stack.Numbers.Small)-1])
		stack.Numbers.Small = stack.Numbers.Small[:len(stack.Numbers.Small)-1]
	} else {
		stack.Numbers.Grow()
		number = stack.Numbers.Big[len(stack.Numbers.Big)-1]
		stack.Numbers.Big = stack.Numbers.Big[:len(stack.Numbers.Big)-1]
	}
	return
}

func (stack *Stack) Heap() {
	address := stack.Pull()
	
	switch {
		case address.Small == 0:
			if len(stack.HeapRoom) > 0 {
				address := stack.HeapRoom[len(stack.HeapRoom)-1]
				stack.HeapRoom = stack.HeapRoom[:len(stack.HeapRoom)-1]
				
				stack.TheHeap[(address%(len(stack.TheHeap)+1))-1] = stack.Grab()
				stack.Push(NewNumber(address))
				
			} else {
				stack.TheHeap = append(stack.TheHeap, stack.Grab())
				stack.Push(NewNumber(len(stack.TheHeap)))
			}
			
		case address.Small > 0:
			stack.Share(stack.TheHeap[(int(address.Small)%(len(stack.TheHeap)+1))-1])
			
		case address.Small < 0:
			stack.TheHeap[(-int(address.Small))%(len(stack.TheHeap)+1)-1] = Array{}
			stack.HeapRoom = append(stack.HeapRoom, -int(address.Small))
	}
}


func (stack *Stack) Put(number Number) {
	var array = (*stack.ActiveArray)
	if number.Int == nil && number.Small < 256  && number.Small >= 0 && array.Big == nil {
		(*stack.ActiveArray).Small = append(array.Small, byte(number.Small))
	} else {
		(*stack.ActiveArray).Grow()
		(*stack.ActiveArray).Big = append(array.Big, number)
	}
}
func (stack *Stack) Pop() (number Number) {
	var array = (*stack.ActiveArray)
	if array.Big == nil {
		number.Small = int64(array.Small[len(array.Small)-1])
		(*stack.ActiveArray).Small = array.Small[:len(array.Small)-1]
	} else {
		array.Grow()
		number = array.Big[len(array.Big)-1]
		(*stack.ActiveArray).Big = array.Big[:len(array.Big)-1]
	}
	return
}

func (stack *Stack) Place(array Array) {
	stack.ActiveArray = &array
}
func (stack *Stack) Get() (number Number) {
	var array = (*stack.ActiveArray)
	var index Number
	
	if array.Big == nil {
		index.Mod(stack.Pull(), NewNumber(len(array.Small)))
		number = NewNumber(int(array.Small[index.ToInt()]))
	} else {
		index.Mod(stack.Pull(), NewNumber(len(array.Big)))
		number = array.Big[index.ToInt()]
	}
	return
}
func (stack *Stack) Set(number Number) {
	var array = (*stack.ActiveArray)
	var index Number
	
	if number.Int == nil && array.Big == nil && number.Small < 256  && number.Small >= 0 {
		index.Mod(stack.Pull(), NewNumber(len((*stack.ActiveArray).Small)))
		(*stack.ActiveArray).Small[index.ToInt()] = byte(number.Small)
	} else {
		(*stack.ActiveArray).Grow()
		index.Mod(stack.Pull(), NewNumber(len((*stack.ActiveArray).Big)))
		(*stack.ActiveArray).Big[index.ToInt()] = number
	}
}

func (stack *Stack) Load() {
	var name string
	var variable string
	var result = stack.Array()
	var err error
	
	text := stack.Grab()
	
	if text.Index(0).ToInt() == '$' && text.Len().Small > 1 {
	
		for i := 0; i < int(text.Len().Small); i++ {
			if i == 0 {
				continue
			}
			name += string(rune(text.Index(i).ToInt()))
		}
		variable = os.Getenv(name)
	} else {
	
		name = text.String()
		
		protocol := strings.SplitN(name, "://", 2)
		if len(protocol) > 1 {
			switch protocol[0] {
				case "tcp":
					listener, err := net.Listen("tcp", ":"+protocol[1])
					_, variable, _ = net.SplitHostPort(listener.Addr().String())
					if protocol[1] == "0" {
						Networks_In[variable] = listener
					} else {
						Networks_In[protocol[1]] = listener
					}
					if err != nil {
						stack.ERROR = NewNumber(1)
					}
				case "dns":
					//This can be optimised. Check the string.
					hosts, err := net.LookupAddr(protocol[1])
					if err != nil {
						hosts, err = net.LookupHost(protocol[1])
						if err != nil {
							stack.ERROR = NewNumber(1)
						}
					}
					variable = strings.Join(hosts, " ")
				default:
					if err != nil {
						stack.ERROR = NewNumber(1)
					}
			}
		} else {
	
			if len(os.Args) > int(text.Index(0).ToInt()) {
				variable = os.Args[text.Index(0).ToInt() ]
			} else {
				stack.ERROR = NewNumber(1)
			}
		}
	}
	
	result.Small = []byte(variable)
	stack.Share(result)
}

func (stack *Stack) Open() {
	var err error
	
	text := stack.Grab()

	var filename string = text.String()
	
	
	var it Pipe
	it.Name = filename
	
	protocol := strings.SplitN(filename, "://", 2)
	if len(protocol) > 1 {
		switch protocol[0] {
		
			case "tcp":
				if listener, ok := Networks_In[protocol[1]]; ok {
					
					it.Connection, err = listener.Accept()
					it.Pipe = it.Connection
					if err != nil {
						stack.Push(NewNumber(-1))
						stack.Relay(it)
						return
					}
					stack.Push(NewNumber(0))
					stack.Relay(it)
					return
					
				} else {
					it.Connection, err = net.Dial("tcp", protocol[1])
					it.Pipe = it.Connection
					if err != nil {
						stack.Push(NewNumber(-1))
						stack.Relay(it)
						return
					}
					stack.Push(NewNumber(0))
					stack.Relay(it)
					return
				}
		}
	}

	it.Pipe, err = os.OpenFile(filename, os.O_RDWR|os.O_APPEND, 0666)
	if err == nil {
		stack.Push(NewNumber(0))
		stack.Relay(it)
		return
	}
	if _, err = os.Stat(filename); err == nil {
		stack.Push(NewNumber(0))
		stack.Relay(it)
		return
	}
	stack.Push(NewNumber(-1))
	stack.Relay(it)
	return
}

func (stack *Stack) Info() {
	var request string
	var variable string
	
	var result = stack.Array()

	text := stack.Grab()
	it := stack.Take()
	
	request = text.String()
	
	switch request {
		case "address":
			if it.Connection != nil {
				variable = it.Connection.RemoteAddr().String()
			}
		case "ip":
			if it.Connection != nil {
				variable = it.Connection.RemoteAddr().(*net.TCPAddr).IP.String()
			}
		case "port":
			if it.Connection != nil {
				variable = strconv.Itoa(it.Connection.RemoteAddr().(*net.TCPAddr).Port)
			}
	}
	
	result.Small = []byte(variable)
	stack.Share(result)
}

func (stack *Stack) Out() {
	var err error
	
	text := stack.Grab()
	f := stack.Take()
	
	if f.Pipe == nil {
		if f.Name[len(f.Name)-1] == '/' {
			i, err := os.Stat(f.Name)
			if err == nil && i.IsDir() {
				
			} else {
				err := os.Mkdir(f.Name, 0666)
				if err != nil {
					stack.Push(NewNumber(-1))
					return
				}
			}
		} else {
			f.Pipe, err = os.Create(f.Name)
			if err != nil {
				stack.Push(NewNumber(-1))
				return
			}
		}
	}
	if int(text.Len().Small) == 0 {
		stack.Push(NewNumber(0))
		return
	}
	for i := 0; i < int(text.Len().Small); i++ {
		v := text.Index(i)
		if v.Int != nil {
			_, err := f.Pipe.Write(v.Bytes())
			if err != nil {
				stack.Push(NewNumber(-1))
				return
			}
		} else {
			_, err := f.Pipe.Write([]byte{byte(v.Small)})
			if err != nil {
				stack.Push(NewNumber(-1))
				return
			}
		}
	} 
	
	stack.Push(NewNumber(0))
}

func (stack *Stack) Stdout() {
	text := stack.Grab()
	
	if text.Big == nil {
		os.Stdout.Write(text.Small)
		return
	}
	
	for i := 0; i < int(text.Len().Small); i++ {
		v := text.Index(i)
		if v.Int != nil {
			os.Stdout.Write(v.Bytes())
		} else {
			os.Stdout.Write([]byte{byte(v.Small)})
		}
	} 
}

func (stack *Stack) In() {

	length := stack.Pull()
	f := stack.Take()
	
	
	if f.Pipe == nil {
		stack.Push(NewNumber(-1000))
		return
	}
	var err error
	var b []byte = make([]byte, int(length.ToInt()))
	var n int
	n, err = f.Pipe.Read(b)
	if len(b) > 1 || n == 0 {
		stack.Push(NewNumber(-1000))
	}
	if err != nil {
		//println(err.Error())
	}
	for i:=n-1; i>=0; i-- {
		stack.Push(NewNumber(int(b[i])))
	}

}

var stdin_reader = bufio.NewReader(os.Stdin)

func  (stack *Stack) Stdin() {
	length := stack.Pull()
	
	switch {
		case length.Small == 0:
			
			b, err := stdin_reader.ReadBytes('\n')
			if err != nil || len(b) == 0 {
				stack.Share(Array{})
			}
			stack.Share(Array{Small:b[:len(b)-1]})
			
		case length.Small > 0:
		
			var b []byte = make([]byte, int(length.ToInt()))
			n, err := os.Stdin.Read(b)
			if err != nil || n <= 0 {
				stack.Share(Array{})
			}
			stack.Share(Array{Small:b})
			
		case length.Small < 0:
			b, err := stdin_reader.ReadBytes(byte(-length.Small))
			if err != nil || len(b) == 0 {
				stack.Share(Array{})
			}
			stack.Share(Array{Small:b[:len(b)-1]})
	}
}

type Number struct {
	*big.Int
	Small int64
}

func (z *Number) Init() {

}

func NewNumber(n int) Number {
	return Number{Small:int64(n)}
}

func (z *Number) Add(a, b Number) {
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Add(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Add(big.NewInt(a.Small), b.Int)
		} else {
			if (a.Small > 0 && b.Small > (1<<63 - 1) - a.Small) {
				z.Int.Add(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			} else if (a.Small < 0 && b.Small < (-1 << 63) - a.Small) {
				z.Int.Add(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			}
			*z = Number{Small:a.Small+b.Small}
		}
	} else {
		z.Int.Add(a.Int, b.Int)
	}
} 

func (z *Number) Sub(a, b Number) {
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Sub(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Sub(big.NewInt(a.Small), b.Int)
		} else {
			if (a.Small > 0 && -b.Small > (1<<63 - 1) - a.Small) {
				z.Int.Sub(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			} else if (a.Small < 0 && -b.Small < (-1 << 63) - a.Small) {
				z.Int.Sub(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			}
			*z = Number{Small:a.Small-b.Small}
		}
	} else {
		z.Int.Sub(a.Int, b.Int)
	}
} 

func (z *Number) Mul(a, b Number) {
	
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Mul(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Mul(big.NewInt(a.Small), b.Int)
		} else {
			*z = Number{Small:a.Small*b.Small}
			if (a.Small != 0 && z.Small / a.Small != b.Small) {
				z.Int = big.NewInt(0)
				z.Int.Mul(big.NewInt(a.Small), big.NewInt(b.Small))
			}
		}
	} else {
		z.Int.Mul(a.Int, b.Int)
	}
} 

var Zero_go = big.NewInt(0)

func (z *Number) Div(a, b Number) {
	defer func() {
        if r := recover(); r != nil {
		    if a.Small == 0 || a.Int != nil && a.Int.Cmp(Zero_go) == 0 {
		    	var b []byte = []byte{1}
		    	rand.Read(b)
		    	*z =NewNumber(int(b[0]+1))
		    	return
		    }
		   	*z = NewNumber(0)
		}
    }()
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Div(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Div(big.NewInt(a.Small), b.Int)
		} else {
			if b.Small < 0 && a.Small == (-1 << 63) {
				z.Int = big.NewInt(0)
				z.Int.Div(big.NewInt(a.Small), big.NewInt(b.Small))
				return
			}
			*z = Number{Small:a.Small/b.Small}
		}
	} else {
		z.Int.Div(a.Int, b.Int)
	}
} 

func (z *Number) Mod(a, b Number) {
	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Mod(a.Int, big.NewInt(b.Small))
		} else if !cb {
			z.Int.Mod(big.NewInt(a.Small), b.Int)
		} else {
			z.Int = big.NewInt(0)
			z.Int.Mod(big.NewInt(a.Small), big.NewInt(b.Small))
			return
		}
	} else {
		z.Int.Mod(a.Int, b.Int)
	}
} 

func (z *Number) Pow(a, b Number) {
	if a.ToInt() == 0 {
		z.Mod(b, NewNumber(2))
		if z.ToInt() != 0 {
			var b []byte = []byte{1}
	    	rand.Read(b)
	    	*z = NewNumber(int(b[0]))
	    	return
		}
		z.Int = big.NewInt(0)
		return
	}

	z.Int = big.NewInt(0)
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			z.Int.Exp(a.Int, big.NewInt(b.Small), nil)
		} else if !cb {
			z.Int.Exp(big.NewInt(a.Small), b.Int, nil)
		} else {
			z.Int = big.NewInt(0)
			z.Int.Exp(big.NewInt(a.Small), big.NewInt(b.Small), nil)
			return
		}
	} else {
		z.Int.Exp(a.Int, b.Int, nil)
	}
} 

func (a Number) Slt(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) == -1 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) == -1 ))
		} else {
			return NewNumber(__bool2int( a.Small < b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) == -1 ))
	}
}

func (a Number) Seq(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) == 0 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) == 0 ))
		} else {
			return NewNumber(__bool2int( a.Small == b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) == 0 ))
	}
}

func (a Number) Sge(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) >= 0 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) >= 0 ))
		} else {
			return NewNumber(__bool2int( a.Small >= b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) >= 0 ))
	}
}

func (a Number) Sle(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) <= 0 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) <= 0 ))
		} else {
			return NewNumber(__bool2int( a.Small <= b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) <= 0 ))
	}
}

func (a Number) Sgt(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) == 1 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) == 1 ))
		} else {
			return NewNumber(__bool2int( a.Small > b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) == 1 ))
	}
}

func (a Number) Sne(b Number) Number {
	if ca, cb := a.Int == nil, b.Int == nil; ca || cb {
		if !ca {
			return NewNumber(__bool2int( a.Cmp(big.NewInt(b.Small)) != 0 ))
		} else if !cb {
			return NewNumber(__bool2int( big.NewInt(a.Small).Cmp(b.Int) != 0 ))
		} else {
			return NewNumber(__bool2int( a.Small != b.Small ))
		}
	} else {
		return NewNumber(__bool2int( a.Cmp(b.Int) != 0 ))
	}
}

func (z Number) True() bool {
	if z.Int == nil {
		if z.Small != 0 {
			return true
		}
		return false
	}
	return z.Int.Cmp(big.NewInt(0)) != 0
} 

func (z Number) ToInt() int {
	if z.Int == nil {
		return int(z.Small)
	}
	return int(z.Int.Int64())
} 


type Pipe struct {
	Name string
	
	Pipe io.ReadWriteCloser
	Connection net.Conn
	
	Function func(*Stack)
}

func (z *Pipe) Init() {

}

func (pipe *Pipe) Exe(stack *Stack) {
	pipe.Function(stack)
}

func (pipe *Pipe) Close() {
	if pipe.Pipe != nil {
		pipe.Pipe.Close()
	}
}

type Array struct {
	Big []Number
	Small []byte
}

func (z *Array) Init() {

}

func (z *Array) Grow() {
	if z.Big == nil {
		z.Big = make([]Number, len(z.Small))
		for i := range z.Small {
			z.Big[i] = Number{Small:int64(z.Small[i])}
		}
	}
}

func (z *Array) Index(n int) Number {
	if z.Big == nil {
		return Number{Small:int64(z.Small[n])}
	} else {
		return z.Big[n]
	}
}

func (z *Array) Len() Number {
	if z.Big == nil {
		return Number{Small:int64(len(z.Small))}
	} else {
		return Number{Small:int64(len(z.Big))}
	}
}

func (z *Array) String() string {
	if z.Big == nil {
		return string(z.Small)
	} else {
		var name string
		for i := 0; i < int(z.Len().Small); i++ {
			name += string(rune(z.Big[i].ToInt()))
		}
		return name
	}
}

func NewStringArray(s string) Array {
	var result = Array{Small:[]byte(s)}
	return result
}

func (array *Array) Join(b Array) Array {
	switch {
		case array.Small != nil && b.Small != nil:
			return Array{Small:append(array.Small, b.Small...)}
			
		case (array.Small != nil && b.Big != nil) || (array.Big != nil && b.Small != nil):
			array.Grow()
			b.Grow()
			fallthrough
		case array.Big != nil && b.Big != nil:
			return Array{Big:append(array.Big, b.Big...)}
	}
	return Array{}
}

func __bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}
`
