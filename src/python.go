package uct

import "flag"

var PythonReserved = []string{
	"and",       "del",       "from",      "not",       "while",   
	"as",        "elif",      "global",    "or",        "with",    
	"assert",    "else",      "if",        "pass",      "yield",    
	"break",     "except",    "import",    "print",     "len",          
	"class",     "exec",      "in",        "raise", 	"open",             
	"continue",  "finally",   "is",        "return",    "bool",   
	"def",       "for",       "lambda",    "try",		"copy",
	"list",

}

//This is the Java compiler for uct.
var Python bool

func init() {
	flag.BoolVar(&Python, "py", false, "Target Python")
	
	instruction := PythonAssembly["IF"]
	instruction.Pass = "pass\n"
	PythonAssembly["IF"] = instruction

	instruction = PythonAssembly["ELSE"]
	instruction.Pass = "pass\n"
	instruction.Check = "IF"
	PythonAssembly["ELSE"] = instruction
	
	instruction = PythonAssembly["RETURN"]
	instruction.Pass = "pass\n"
	instruction.Check = "FUNCTION"
	PythonAssembly["RETURN"] = instruction

	RegisterAssembler(PythonAssembly, &Python, "py", "#")

	for _, word := range PythonReserved {
		PythonAssembly[word] = Reserved()
	}
}

var PythonAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   "#! /bin/python3\nimport stack\nimport sys\nfrom multiprocessing import Process\ndef i_div(a, b):\treturn a//b",
		Args: 1,
	},

	"FOOTER": Instruction{},

	"FILE": Instruction{
		Path: "stack.py",
	},
	
	"PYTHON": Instruction{All:true},

	"NUMBER": is("%s", 1),
	"BIG": 	is("%s", 1),
	"SIZE":   is("len(%s)", 1),
	"STRING": is("list(bytes(%s, 'utf-8'))", 1),
	"ERRORS":  is("stack.ERROR", 1),
	
	"LINK":  is("stack.take().data = stack.grab()"),
	"CONNECT":  is("stack.share(stack.take().data)"),
	
	"SLICE":  is("stack.slice()"),

	"SOFTWARE": Instruction{
		Data:   "if __name__ == '__main__':\n\tstack = stack.Stack()\n",
		Indent: 1,
	},
	"EXIT": Instruction{
		Data:        "sys.exit(stack.ERROR)",
	},

	"FUNCTION": is("def %s(stack):", 1, 1),
	"RETURN": Instruction{
		Pass: "pass\n",
		
		Indented:    1,
		Data:        "\n",
		Indent:      -1,
		Else: &Instruction{
			Data: "return",
		},
	},
	
	"SCOPE": is(`stack.relay(stack.pipe(%s))`, 1),
	
	"EXE": is("%s.exe(stack)", 1),

	"PUSH": is("stack.numbers.append(%s)", 1),
	"PULL": is("%s = stack.numbers.pop()", 1),

	"PUT":   is("stack.activearray.append(%s)", 1),
	"POP":   is("%s = stack.activearray.pop()", 1),
	"PLACE": is("stack.activearray = %s", 1),
	"ARRAY":  is("%s = stack.array()", 1),
	"RENAME": is("%s = stack.grab()", 1),
	"RELOAD": is("%s = stack.take()", 1),

	"SHARE": is("stack.arrays.append(%s)", 1),
	"GRAB":  is("%s = stack.arrays.pop()", 1),

	"RELAY": is("stack.pipes.append(%s)", 1),
	"TAKE":  is("%s = stack.pipes.pop()", 1),

	"GET": is("%s = stack.get()", 1),
	"SET": is("stack.set(%s)", 1),

	"VAR": is("%s = 0", 1),

	"PIPE": is("%s = stack.queue()", 1),
	
	"INBOX":   is("stack.share(stack.inbox.recv())\n"),
	"READYBOX": is("stack.push(stack.inbox.poll())"),
	"OUTBOX":   is("stack.outbox.send(stack.grab())"),

	"EVAL": is("eval(bytes(stack.grab()).decode()+'(stack)')"),

	"OPEN":   is("stack.open()"),
	"EXECUTE": is("stack.execute()"),
	"DELETE": is("stack.delete()"),
	"LOAD":   is("stack.load()"),
	"OUT":    is("stack.out()"),
	"STAT":   is("stack.info()"),
	"IN":     is("stack.inn()"),
	"STDOUT": is("stack.stdout()"),
	"STDIN":  is("stack.stdin()"),
	"HEAP":   is("stack.heap()"),
	"HEAPIT":   is("stack.heapit()"),
	"MAKE":   is("stack.share([0]*stack.pull())"),

	"CLOSE": is("%s.close()", 1),

	"LOOP":   is("while 1:", 0, 1),
	"BREAK":  is("break"),
	"REPEAT": is("\n", 0, -1, -1),

	"IF":   is("if %s != 0:", 1, 1),
	"ELSE": is("else:", 0, 0, -1),
	"END":  is("\n", 0, -1, -1),

	"RUN":  is("%s(stack)", 1),
	"DATA": is_data("%s = [", "%s", ",", "]"),

	"FORK": is("Process(target=%s, args=(stack.copy(),)).start()\n", 1),

	"ADD": is("%s = %s + %s", 3),
	"SUB": is("%s = %s - %s", 3),
	"MUL": is("%s = stack.mul(%s, %s)", 3),
	"DIV": is("%s = stack.div(%s, %s)", 3),
	"MOD": is("%s = stack.mod(%s, %s)", 3),
	"POW": is("%s = stack.pow(%s, %s)", 3),

	"SLT": is("%s = int(%s <  %s)", 3),
	"SEQ": is("%s = int(%s == %s)", 3),
	"SGE": is("%s = int(%s >= %s)", 3),
	"SGT": is("%s = int(%s >  %s)", 3),
	"SNE": is("%s = int(%s != %s)", 3),
	"SLE": is("%s = int(%s <= %s)", 3),

	"JOIN": is("%s = %s + %s", 3),
	"ERROR": is("stack.ERROR = %s", 1),
}
