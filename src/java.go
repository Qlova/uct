package main

import "flag"
import "fmt"

var JavaReserved = []string{
	"abstract", "continue", "for", "new", "switch", "assert", "default",
	"goto", "package", "synchronized", "boolean", "do", "if", "private",
	"this", "break", "double", "implements", "protected", "throw",
	"byte", "else", "import", "public", "throws", "case", "enum", "instanceof",
	"return", "transient", "catch", "extends", "int", "short", "try",
	"char", "final", "interface", "static", "void", "class", "finally",
	"long", "strictfp", "volatile", "const", "float", "native", "super", "while",
}

//This is the Java compiler for uct.
var Java bool

func init() {
	flag.BoolVar(&Java, "java", false, "Target Java")

	RegisterAssembler(JavaAssembly, &Java, "java", "//")

	for _, word := range JavaReserved {
		JavaAssembly[word] = Reserved()
	}
}

var JavaAssembly = Assemblable{
	//Special commands.
	"HEADER": Instruction{
		Data:   "public class %s {",
		Indent: 1,
		Args:   1,
	},

	"FOOTER": Instruction{
		Data:        "}",
		Indent:      -1,
		Indentation: -1,
	},

	"FILE": Instruction{
		Data: JavaFile,
		Path: "/Stack.java",
	},

	"NUMBER": is("new Stack.Number(%s)", 1),
	"BIG": is("new Stack.Number(%s)", 1),
	"SIZE":   is("%s.size()", 1),
	"STRING": is("new Stack.Array(%s)", 1),
	"ERRORS":  is("stack.ERROR", 1),

	"SOFTWARE": Instruction{
		Data:   "public static void main(String[] args) { Stack stack = new Stack(); stack.Arguments = args;",
		Indent: 1,
	},
	"EXIT": Instruction{
		Indented:    2,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "System.exit(stack.ERROR.intValue());",
		},
	},

	"FUNCTION": is("static void %s(Stack stack) {", 1, 1),
	"RETURN": Instruction{
		Indented:    2,
		Data:        "}\n",
		Indent:      -1,
		Indentation: -1,
		Else: &Instruction{
			Data: "return;",
		},
	},
	
	"SCOPE": is(`Class[] cArg = new Class[1]; cArg[0] = Stack.class; try { stack.relay(new Stack.Pipe((new Object() { }.getClass().getEnclosingClass().getDeclaredMethod("%s", cArg)))); } catch (NoSuchMethodException e) { throw new RuntimeException(e); }`, 1),
	
	"EXE": is("%s.exe(stack);", 1),

	"PUSH": is("stack.push(%s);", 1),
	"PULL": is("Stack.Number %s = stack.pull();", 1),

	"PUT":   is("stack.put(%s);", 1),
	"POP":   is("Stack.Number %s = stack.pop();", 1),
	"PLACE": is("stack.place(%s);", 1),

	"ARRAY":  is("Stack.Array %s = stack.array();", 1),
	"RENAME": is("%s = stack.ActiveArray;", 1),
	
	"RELOAD": is("%s = stack.take();", 1),

	"SHARE": is("stack.share(%s);", 1),
	"GRAB":  is("Stack.Array %s = stack.grab();", 1),

	"RELAY": is("stack.relay(%s);", 1),
	"TAKE":  is("Stack.Pipe %s = stack.take();", 1),

	"GET": is("Stack.Number %s = stack.get();", 1),
	"SET": is("stack.set(%s);", 1),

	"VAR": is("Stack.Number %s = new Stack.Number();", 1),

	"OPEN":   is("stack.open();"),
	"LOAD":   is("stack.load();"),
	"OUT":    is("stack.out();"),
	"STAT":   is("stack.info();"),
	"IN":     is("stack.in();"),
	"STDOUT": is("stack.stdout();"),
	"STDIN":  is("stack.stdin();"),
	"HEAP":   is("stack.heap();"),
	"LINK":   is("stack.link();"),
	"CONNECT":   is("stack.connect();"),
	"SLICE":   is("stack.slice();"),

	"CLOSE": is("%s.close();", 1),

	"LOOP":   is("while (true) {", 0, 1),
	"BREAK":  is("break;"),
	"REPEAT": is("}", 0, -1, -1),

	"IF":   is("if (%s.compareTo(new Stack.Number(0)) != 0 ) {", 1, 1),
	"ELSE": is("} else {", 0, 0, -1),
	"END":  is("}", 0, -1, -1),

	"RUN":  is("%s(stack);", 1),
	"DATA": is("static Stack.Array %s = %s;", 2),

	"FORK": is("{ Stack s = stack.copy(); Stack.ThreadPool.execute(() -> %s(s)); }\n", 1),

	"ADD": is("%s = %s.add(%s);", 3),
	"SUB": is("%s = %s.sub(%s);", 3),
	"MUL": is("%s = %s.mul(%s);", 3),
	"DIV": is("%s = %s.div(%s);", 3),
	"MOD": is("%s = %s.mod(%s);", 3),
	"POW": is("%s = %s.pow(%s);", 3),

	"SLT": is("%s = %s.compareTo(%s) == -1 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SEQ": is("%s = %s.compareTo(%s) == 0 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SGE": is("%s = %s.compareTo(%s) >= 0 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SGT": is("%s = %s.compareTo(%s) == 1 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SNE": is("%s = %s.compareTo(%s) != 0 ? new Stack.Number(1) : new Stack.Number(0);", 3),
	"SLE": is("%s = %s.compareTo(%s) <= 0 ? new Stack.Number(1) : new Stack.Number(0);", 3),

	"JOIN": is("%s = %s.join(%s);", 3),
	"ERROR": is("stack.ERROR = %s;", 1),
}

func JavaScope(args []string) string {
		filename := JavaAssembly["NAME"].Data
		return fmt.Sprintf(`Class[] cArg = new Class[1]; cArg[0] = Stack.class; try { stack.relay(new Pipe(new Object() { }.getClass().getEnclosingClass().getDeclaredMethod("%s", cArg))); } catch (NoSuchMethodException e) { throw new RuntimeException(e); }`, filename, args[0])
}

//Edit this in a Java IDE.
const JavaFile = `

//Compiled to Java with UCT (Universal Code Translator)

//Import java libraries.
import java.math.BigInteger; 	//Support numbers of any size.
import java.util.Hashtable;  	//Hashtable is a useful utilPipey.
import java.util.ArrayList;  	//ArrayLists are helpful.
import java.io.*;				//Deal with files.
import java.net.*;				//Deal with the network.
import java.lang.reflect.*;		//Reflection for methods.

//Random numbers.
import java.security.SecureRandom;

//This is for threading.
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;

//This is the Java stack implementation.
// It holds arrays for the 4 types:
//		Numbers
//		Arrays
//		Functions
//		Pipes
//
// It also holds the ERROR variable for the current thread.
// The currently active array is stored as ActiveArray.
public class Stack {
    Array 				Numbers;
    ArrayArray 			Arrays;
    PipeArray 			Pipes;

    Number 				ERROR;

    Array				ActiveArray;
    ArrayArray			Heap;
    ArrayList<Integer>	HeapRoom;
    

    //This hashtable keeps track of Servers currently listening on the specified port.
    static Hashtable<String, ServerSocket> Networks_In = new Hashtable<String, ServerSocket>();

    //This will store the system arguments.
    public static String[] Arguments;

    //This is the threading pool.
    public static Executor ThreadPool = Executors.newCachedThreadPool();

    //This creates an empty stack. Ready for use.
    public Stack() {
        Numbers 		= new Array();
        Arrays 			= new ArrayArray();
        Pipes 		    = new PipeArray();
        ERROR 			= new Number();

        Heap			= new ArrayArray();
        HeapRoom		= new ArrayList<Integer>();
    }

    //This returns a copy of a stack which can be used by another thread.
    public Stack copy() {
        Stack n = new Stack();
        n.Numbers = new Array();

        n.Numbers.List = new ArrayList<Number>(Numbers.List);

        n.Arrays.List = new ArrayList<Array>(Arrays.List.size());
        for (int i = 0; i < Arrays.List.size()-1; i++) {
            n.Arrays.List.set(i, new Array());
            n.Arrays.List.get(i).List = new ArrayList<>(Arrays.List.get(i).List);
        }

        n.Pipes = new PipeArray();
        n.Pipes.List = new ArrayList<>(Pipes.List);
        return n;
    }

    //This stuff can be inlined in the compiler, they are only here for reference.

    //SHARE array
    void share(Array a) {
        Arrays.push(a);
    }

    //GRAB array
    Array grab() {
        return Arrays.pop();
    }

    //RELAY pipe
    void relay(Pipe p) {
        Pipes.push(p);
    }

    //TAKE pipe
    Pipe take() {
        return Pipes.pop();
    }

    //PUSH number
    void push(Number n) {
        Numbers.push(n);
    }

    //PUT number
    void put(Number b) {
        ActiveArray.push(b);
    }

    //PLACE array
    void place(Array a) {
        ActiveArray = a;
    }

    //ARRAY name
    Array array() {
        Array a = new Array();
        ActiveArray = a;
        return a;
    }

    //GET number
    Number get() {
        return ActiveArray.index(pull().mod(ActiveArray.size()));
    }

    //SET number
    void set(Number b) {
        ActiveArray.set(pull(), b);
    }

    //POP number
    Number pop() {
        return ActiveArray.pop();
    }

    //PULL number
    Number pull() {
        return Numbers.pop();
    }
    
    void heap() {
    	Number address = pull();
	
		if (address.intValue() == 0) {
			if (HeapRoom.size() > 0) {
				Integer address2 = HeapRoom.remove(HeapRoom.size()-1);
				
				Heap.List.set(((address2)%(Heap.List.size()+1)-1),  grab());
				push(new Number(address2));
			} else {
				Heap.push(grab());
				push(Heap.size());
			}
			
		} else if (address.intValue() > 0) {
			share(Heap.List.get((address.intValue()%(Heap.List.size()+1))-1));
			
		} else if (address.intValue() < 0) {
			Heap.List.set(((-address.intValue())%(Heap.List.size()+1)-1),  null);
			HeapRoom.add(-address.intValue());
		}
    }

    void stdout() {
        Array text = grab();
        for (int i = 0; i < text.size().intValue(); i++) {
            if (text.index(new Number(i)) != null) {
                int c = text.index(new Number(i)).intValue();
                System.out.print((char)(c));
            }
        }
    }

    void stdin() {
        Number length = pull();
        
        //This is the mode we use.
        // >0 is number of bytes to read.
        // <0 is character to read.
        // 0 is read a line.
        if (length.compareTo(new Number(0)) == 0) {
        
        	byte[] b = new byte[1];
        	
        	String input = "";
        	
        	while (true) {
        		try {
        			int n = System.in.read(b);
        		} catch (Exception e) {
        			ERROR = new Number(-1);
        			share(new Array(input));
        			return;
        		}
        		if (b[0] == '\n') {
        			share(new Array(input));
        			return;
        		}
        		input += ((char)(b[0]));
        	}
        	
        	
        } else if (length.compareTo(new Number(0)) == -1)  {
        
       		byte[] b = new byte[1];
        	
        	String input = "";
        	
        	while (true) {
        		try {
        			int n = System.in.read(b);
        		} catch (Exception e) {
        			ERROR = new Number(-1);
        			share(new Array(input));
        			return;
        		}
        		if (b[0] == -length.intValue()) {
        			share(new Array(input));
        			return;
        		}
        		input += ((char)(b[0]));
        	}
        
        
        } else { //length is > 0
        	
        	
        	byte[] b = new byte[length.intValue()];
        	
    		try {
    			int n = System.in.read(b);
    		} catch (Exception e) {
    			ERROR = new Number(-1);
    			share(new Array(b));
    			return;
    		}
        	
			 share(new Array(b));
			 return;
        }
       
    }

    void in() {
        Pipe file = take();
        Number length = pull();
        int n = 0;
        byte[] b = new byte[length.intValue()];

        if (file.input != null) {
            try {
                n = file.input.read(b);
            }catch(Exception e){
                push(new Number(-1000));
            }
        }

        if ((b.length > 1) || (n <= 0)) {
            push(new Number(-1000));
        }

        for (int i = n-1; i >= 0; i--) {
            push(new Number(b[i]));
        }

        return;
    }

    void out() {
        Pipe file = take();
        Array text = grab();

        if (file.output != null) {
            //TODO optimise to send in a single packet.
            for (int i = 0; i < text.size().intValue(); i++) {
                if (text.index(new Number(i)) != null) {
                    int c = text.index(new Number(i)).intValue();
                    try {
                        file.output.write((char)(c));
                    }catch(Exception e){
                        push(new Number(-1));
                    }
                }
            }
            push(new Number(0));
            return;
        }

        if (text.size().intValue() == 0 || file.output == null ) {
            if (file.Name.charAt(file.Name.length()-1) == '/') {
                if (new File(file.Name).exists()) {

                } else {
                    try {
                        File f = new File(file.Name);
                        if (!f.mkdir()) {
                            push(new Number(-1));
                            return;
                        }
                        push(new Number(0));
                        return;
                    } catch (Exception e) {
                        push(new Number(-1));
                        return;
                    }
                }
            } else if (new File(file.Name).exists()) {

            } else {
                try {
                    new File(file.Name).createNewFile();
                    push(new Number(0));
                    return;
                } catch (Exception e)  {
                    push(new Number(-1));
                    return;
                }
            }
        }

        for (int i = 0; i < text.size().intValue(); i++) {
            if (text.index(new Number(i)) != null) {
                int c = text.index(new Number(i)).intValue();
                try {
                    file.output.write((char)(c));
                }catch(Exception e){
                    push(new Number(-1));
                }
            }
        }
        push(new Number(0));
    }


    //LOAD
    void load() {
        String name = "";
        String variable = "";

        Array result = new Array();

        //This request is what we need to load.
        Array request = grab();

        //The request is an enviromental variable.
        if (request.index(0) == '$' && request.length() > 1) {

            //We parse the rest of the string.
            for (int i = 1; i < request.length(); i++) {
                try {
                    name += (char)(request.index(i));
                } catch (Exception e) {

                }
            }
            share(new Array(System.getenv(name)));
            return;
        }


        name = request.String();

        //Load various protocols.
        //Protocols are seperated by ://
        //For example http://, dns://, tcp://
        String[] protocol = name.split("://", 2);
        if (protocol.length > 1) {
            switch (protocol[0]) {
                case "tcp":
                    try {
                        ServerSocket ss = new ServerSocket(Integer.parseInt(protocol[1]));
                        String port = String.valueOf(ss.getLocalPort());
                        if (protocol[1].equals("0")) {
                            Networks_In.put(port, ss);
                            variable = String.valueOf(port);
                        } else {
                            Networks_In.put(protocol[1], ss);
                            variable = protocol[1];
                        }
                        share(new Array(variable));
                        return;
                    } catch (Exception e) {
                        ERROR = new Number(-1);
                    }
                    break;

                case "dns":
                    try {
                        InetAddress address = InetAddress.getByName(protocol[1]);
                        variable = address.getHostName();
                        if (variable == protocol[1]) {
                            variable = "";
                            throw null;
                        }
                        share(new Array(variable));
                        return;
                    } catch (Exception e) {
                        try {
                            InetAddress[] addresses = InetAddress.getAllByName(protocol[1]);
                            for (int i=0; i<addresses.length-1; i++) {
                                variable += addresses[i].getHostAddress();
                                variable += " ";
                            }
                            share(new Array(variable));
                        return;
                        } catch (Exception e2) {
                            ERROR = new Number(-1);
                        }
                    }
            }
        }

        if (Arguments.length > request.index(new Number(0)).intValue()-1) {
            share(new Array(Arguments[request.index(new Number(0)).intValue()-1]));
            return;
        } else {
        	ERROR = new Number(1);
        	share(result);
        	return;
        }

        share(result);
        ERROR = new Number(404);
        return;
    }


    void open() {

        Array request = grab();
        String path = request.String();

        Pipe pipe = new Pipe();
        pipe.Name = path;

        //Load various protocols.
        String[] protocol = pipe.Name.split("://", 2);
        if (protocol.length > 1) {
            switch (protocol[0]) {
                case "tcp":
                    ServerSocket server = Networks_In.get(protocol[1]);
                    if (server != null) {
                        try {
                            Socket client = server.accept();
                            pipe.socket = client;
                            pipe.input = client.getInputStream();
                            pipe.output = client.getOutputStream();
                        } catch( Exception e) {
                            push(new Number(-1));
                            relay(pipe);
                            return;
                        }
                        push(new Number(0));
                         relay(pipe);
                         return;
                    } else {
                        String[] hostport = protocol[1].split(":", 2);
                        if (hostport.length > 1) {
                            //TODO
                            try {
                                Socket req = new Socket(hostport[0], (int)Integer.valueOf(hostport[1]));
                                pipe.socket = req;
                                pipe.input = req.getInputStream();
                                pipe.output = req.getOutputStream();
                            } catch (Exception e) {
                                push(new Number(-1));
                                relay(pipe);
                                return;
                            }
                            push(new Number(0));
                           	relay(pipe);
                           	return;
                        }
                    }
            }
        }

        File file = new File(path);

        if (file.exists()) {
            try {
                pipe.input = new FileInputStream(file);
                pipe.output = new FileOutputStream(file, true);
            }catch(Exception e){
            }
            push(new Number(0));
            relay(pipe);
            return;
        }
        push(new Number(-1));
        relay(pipe);
    }

    void info () {
        String request = "";
        String variable = "";
        Array result = new Array();

        Pipe file = take();
        Array text = grab();

        for (int i = 0; i < text.size().intValue(); i++) {
            if (text.index(new Number(i)) != null) {
                int c = text.index(new Number(i)).intValue();
                request += (char)(c);
            }
        }

        switch (request) {
            case "ip":
                if (file.socket != null) {
                    variable = file.socket.getLocalAddress().getHostAddress();
                }
        }

        for (int i = 0; i < variable.length(); i++) {
            result.push(new Number(variable.charAt(i)));
        }
        share(result);
    }

    public static class Number {
        BigInteger a;

        public Number() {
            a = BigInteger.ZERO;
        }

        public Number(int n) {
            a = BigInteger.valueOf(n);
        }

        public Number(BigInteger n) {
            a = n;
        }
        
        public Number(String s) {
            a = new BigInteger(s);
        }

        public int compareTo(Number b) {
            return a.compareTo(b.a);
        }

        public int intValue() {
            return a.intValue();
        }

        //These functions can be inlined. Only here for reference.
        Number slt(Number b) {
            return new Number(compareTo(b) == -1 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number seq(Number b) {
            return new Number(compareTo(b) == 0 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number sge(Number b) {
            return new Number(compareTo(b) >= 0 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number sgt(Number b) {
            return new Number(compareTo(b) == 1 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number sne(Number b) {
            return new Number(compareTo(b) != 0 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number sle(Number b) {
            return new Number(compareTo(b) <= 0 ? BigInteger.ONE : BigInteger.ZERO);
        }

        Number add(Number b) {
            return new Number(a.add(b.a));
        }

        Number sub(Number b) {
            return new Number(a.subtract(b.a));
        }

        Number mod(Number b) {
            return new Number(a.mod(b.a));
        }

        Number div(Number b) {
            try {
                return new Number(a.divide(b.a));
            } catch (Exception e) {
                if (compareTo(new Number(0)) == 0) {
                    SecureRandom srand = new SecureRandom();
                    return new Number(srand.nextInt(255)+1);
                } else {
                    return new Number(BigInteger.ZERO);
                }
            }
        }

        Number mul(Number b) {
            return new Number(a.multiply(b.a));
        }

        Number pow(Number b) {
            if (intValue() == 0) {
                if (b.mod(new Number(2)).intValue() != 0) {
                    SecureRandom srand = new SecureRandom();
                    return new Number(srand.nextInt(255)+1);
                }
                return new Number(BigInteger.ZERO);
            }
            return new Number(a.pow(b.intValue()));
        }
    }

    //Pipe implementation.
    public static class Pipe {
        String Name;

        //Input and Output.
        InputStream input;
        OutputStream output;

        //Types of pipes.
        File file;
        Socket socket;
        
        Method method;
        
        public Pipe() {
        }
        
        
        public Pipe(Method m) {
        	method = m;
        }
        
        void exe(Stack stack) {
        	try {
        		method.invoke(null, stack); 
        	} catch (Exception e) {  
        		throw new RuntimeException(e);
        	}
        }

        void close() {
            try {
                input.close();
                output.close();
            }catch (Exception e)  {
            }
            try {
                socket.close();
            }catch (Exception e)  {
            }
        }
    }

    public static class PipeArray {
        ArrayList<Pipe> List;

        public PipeArray() {
            List = new ArrayList<Pipe>();
        }

        public Pipe pop() {
            return List.remove(List.size()-1);
        }
        public void push(Pipe n) {
            List.add(n);
        }
    }

    //An array of arrays, or a "Heap".
    public static class ArrayArray {
        ArrayList<Array> List;

        void push(Array n) {
            List.add(n);
        }

        public ArrayArray() {
            List = new ArrayList<Array>();
        }

        Array pop() {
            return List.remove(List.size()-1);
        }


        Number size() {
            return new Number(List.size());
        }

        Array index(Number n) {
            return List.get(n.intValue());
        }
    }

    public static class Array {
        ArrayList<Number> List;
        //new ArrayList<Integer>();

        void push(Number n) {
            List.add(n);
        }

        void set(Number index, Number n) {
            List.set(index.intValue(), n);
        }

        Array join(Array s) {
            Array newList = new Array();
            newList.List.addAll(List);
            newList.List.addAll(s.List);
            return newList;
        }

        public Array(Number... n) {
            List = new ArrayList<Number>();
            for (int i = 0; i < n.length; ++i) {
                List.add(n[i]);
            }
        }

        public Array(String s) {
            List = new ArrayList<Number>();
            for (int i = 0; i < s.length(); ++i) {
                List.add(new Number(s.charAt(i)));
            }
        }
        
        public Array(byte[] s) {
            List = new ArrayList<Number>();
            for (int i = 0; i < s.length; ++i) {
                List.add(new Number(s[i]));
            }
        }

        public String String() {
            String name = "";
            for (int i = 0; i < size().intValue(); i++) {
                if (index(new Number(i)) != null) {
                    int c = index(new Number(i)).intValue();
                    name += (char)(c);
                }
            }
            return name;
        }

        Number pop() {
            return List.remove(List.size()-1);
        }

        Number size() {
            return new Number(List.size());
        }

        Number index(Number n) {
            return List.get(n.intValue());
        }

        int index(int n) {
            return List.get(n).intValue();
        }
        int length() {
            return List.size();
        }
    }
}
`
