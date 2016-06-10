package main

import "errors"
import "flag"
import "strconv"
import "strings"

//Go flag.
var Python bool
func init() {
	flag.BoolVar(&Python, "py", false, "Target Python")
	
	RegisterAssembler(new(PythonAssembler), &Python, "py", "#")
}

type PythonAssembler struct {
	Indentation int
}


func (g *PythonAssembler) indt(n ...int) string {
	if len(n) > 0 {
		return strings.Repeat("\t", int(g.Indentation)+n[0])
	} else {
		return strings.Repeat("\t", int(g.Indentation))
	}
}

func (g *PythonAssembler) Header() []byte {
	return []byte(
	`
import sys
import os

	
N = []
N2 = []
F = []
F2 = []

ERROR = 0

def push(n):
	N.append(n)
	
def pushit(n):
	F2.append(n)

def pushstring(n):
	N2.append(n)

def pushfunc(n):
	F.append(n)

def pop():
	return N.pop()

def popit():
	return F2.pop()

def popstring():
	return N2.pop()
	
def popfunc():
	return F.pop()

def openit():
	filename = ""
	text = popstring()
	for i in range(0, len(text)):
		filename += chr(text[i])
	
	file = None
	try:
		file = open(filename)
	except:
		if os.path.isdir(filename):
			push(0)
		else:
			push(-1)
	else:
		push(0)
	return file

def out(file):
	text = popstring()
	for i in range(0, len(text)):
		file.write(chr(text[i]))

def inn(file):
	length = pop()
	for i in range(0, length):
		v = file.read(1)
		if v == '':
			push(-1000)
			return
		push(ord(v))

def close(file):
	file.close()

def stdout():
	text = popstring()
	for i in range(0, len(text)):
		print(chr(text[i]), end="")

def stdin():
	length = pop()
	for i in range(0, length):
		push(ord(sys.stdin.read(1)))
		
def seq(a, b):
	if a == b:
		return 1
	return 0

def sge(a, b):
	if a >= b:
		return 1
	return 0

def sgt(a, b):
	if a > b:
		return 1
	return 0
	
def sne(a, b):
	if a != b:
		return 1
	return 0

def sle(a, b):
	if a <= b:
		return 1
	return 0

def slt(a, b):
	if a < b:
		return 1
	return 0
`)
}

func (g *PythonAssembler) SetFileName(s string) {
}

func (g *PythonAssembler) Assemble(command string, args []string) ([]byte, error) {
	//Format argument expressions.
	//var expression bool
	for i, arg := range args {
		if arg == "=" {
			//expression = true
			continue
		}
		if _, err := strconv.Atoi(arg); err == nil {
			args[i] = arg
			continue
		}
		if arg[0] == '#' {
			args[i] = "len("+arg[1:]+")"
		}
		if arg[0] == '"' {
			var newarg string
			var j = i
			arg = arg[1:]
			
			stringloop:
			arg = strings.Replace(arg, "\\n", "\n", -1)
			for _, v := range arg {
				if v == '"' {
					goto end
				}
				newarg += strconv.Itoa(int(v))+","
			}
			if len(arg) == 0 {
				goto end
			}
			newarg += strconv.Itoa(int(' '))+","
			j++
			//println(arg)
			arg = args[j]
			goto stringloop
			end:
			//println(newarg)
			args[i] = newarg
		}
		//RESERVED names in the language.
		switch arg {
			case "byte", "len", "open", "close":
				args[i] = "u_"+args[i]
		}
	} 
	switch command {
		case "ROUTINE":
			return []byte(""), nil
		case "SUBROUTINE":
			defer func() { g.Indentation++ }()
			return []byte("def "+args[0]+"():\n"+g.indt(1)+"global ERROR\n"), nil
		case "FUNC":
			return []byte(g.indt()+args[0]+" = "+args[1]+" \n"), nil
		case "EXE":
			return []byte(g.indt()+args[0]+"() \n"), nil
		case "PUSH", "PUSHSTRING", "PUSHFUNC", "PUSHIT":
			var name string
			if command == "PUSHSTRING" {
				name = "string"
			}
			if command == "PUSHFUNC" {
				name = "func"
			}
			if command == "PUSHIT" {
				name = "it"
			}
			if len(args) == 1 {
				return []byte(g.indt()+"push"+name+"("+args[0]+")\n"), nil
			} else {
				return []byte(g.indt()+args[1]+".append("+args[0]+")\n"), nil
			}
		case "POP", "POPSTRING", "POPFUNC", "POPIT":
			var name string
			if command == "POPSTRING" {
				name = "string"
			}
			if command == "POPFUNC" {
				name = "func"
			}
			if command == "POPIT" {
				name = "it"
			}
			if len(args) == 1 {
				return []byte(g.indt()+args[0]+" = pop"+name+"()\n"), nil
			} else {
				return []byte(g.indt()+args[0]+" = "+args[1]+".pop()\n"), nil
			}
		case "INDEX":
			return []byte(g.indt()+args[2]+" = "+args[0]+"["+args[1]+"]\n"), nil
		case "SET":
			return []byte(g.indt()+args[0]+"["+args[1]+"] = "+args[2]+";\n"), nil
		case "VAR":
			if len(args) == 1 {
				return []byte(g.indt()+args[0]+" = 0 \n"), nil
			} else {
				return []byte(g.indt()+args[0]+" = "+args[1]+" \n"), nil
			}
			
		//IT stuff.
		case "OPEN":
			return []byte(g.indt()+""+args[0]+" = openit()\n"), nil
		case "OUT":
			return []byte(g.indt()+"out("+args[0]+")\n"), nil
		case "IN":
			return []byte(g.indt()+"inn("+args[0]+")\n"), nil
		case "CLOSE":
			return []byte(g.indt()+"close("+args[0]+")\n"), nil
		
		case "ERROR":
			return []byte(g.indt()+"ERROR="+args[0]+"\n"), nil
	
		case "STRING":
			return []byte(g.indt()+args[0]+" = [] \n"), nil
		case "STDOUT":
			return []byte(g.indt()+"stdout()\n"), nil
		case "STDIN":
			return []byte(g.indt()+"stdin()\n"), nil
		case "LOOP":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"while True:\n"), nil
		case "REPEAT", "END", "DONE":
			if g.Indentation > 0 {
				g.Indentation--
			}
			return []byte(g.indt()+"\n"), nil
		case "IF":
			defer func() { g.Indentation++ }()
			return []byte(g.indt()+"if "+args[0]+" != 0:\n"), nil
		case "RUN":
			return []byte(g.indt()+args[0]+"()\n"), nil
		case "ELSE":
			return []byte(g.indt(-1)+"else:\n"), nil
		case "ELSEIF":
			return []byte(g.indt(-1)+"elif "+args[0]+" != 0:\n"), nil
		case "BREAK":
			return []byte(g.indt()+"break\n"), nil
		case "RETURN":
			return []byte(g.indt()+"return\n"), nil
		case "STRINGDATA":
			return []byte(g.indt()+args[0]+" = ["+args[1]+"] \n"), nil
		case "JOIN":
			return []byte(g.indt()+args[0]+" = "+args[1]+" + "+args[2]+"\n"), nil
		
		//Maths.
		case "ADD":
			return []byte(g.indt()+args[0]+" = "+args[1]+" + "+args[2]+"\n"), nil
		case "SUB":
			return []byte(g.indt()+args[0]+" = "+args[1]+" - "+args[2]+"\n"), nil
		case "MUL":
			return []byte(g.indt()+args[0]+" = "+args[1]+" * "+args[2]+"\n"), nil
		case "DIV":
			return []byte(g.indt()+args[0]+" = "+args[1]+" // "+args[2]+"\n"), nil
		case "MOD":
			return []byte(g.indt()+args[0]+" = "+args[1]+" % "+args[2]+"\n"), nil
		case "POW":
			return []byte(g.indt()+args[0]+" = "+args[1]+" ** "+args[2]+"\n"), nil
			
		case "SLT", "SEQ", "SGE", "SGT", "SNE", "SLE":
			return []byte(g.indt()+args[0]+" = "+strings.ToLower(command)+"("+args[1]+","+args[2]+")\n"), nil
		default:
			/*if expression {
				return []byte(g.indt()+command+" "+strings.Join(args, " ")+"\n"), nil
			}*/
	}
	return nil, errors.New("Unrecognised expression: "+command+" "+strings.Join(args, " "))

}

func (g *PythonAssembler) Footer() []byte {
	return []byte("")
}
