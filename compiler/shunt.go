package compiler

type Shunt map[string]func(*Compiler, Type) Type

func (c *Compiler) Shunt(t Type, precedence int) Type {
	for peek := c.Peek(); c.GetOperator(peek).Precedence >= precedence; {

		if c.GetOperator(c.Peek()).Precedence == -1 {
			break
		}
		op := c.GetOperator(peek)
		c.Scan()

		rhs := c.Expression()
		peek = c.Peek()
		for c.GetOperator(peek).Precedence > op.Precedence {
			rhs = c.Shunt(rhs, c.GetOperator(peek).Precedence)
			peek = c.Peek()
		}
		
		if t.Shunt != nil {
			if result := t.Shunt(c, op.Symbol, t, rhs); result != nil {
				t = *result
				continue
			}
		}
		
		if shunt, ok := t.Shunts[op.Symbol]; ok {
			t = shunt(c, rhs)
		} else {
			c.RaiseError(Translatable{
				English: "Operator "+op.Symbol+" does not apply to "+t.String(), 
			})
		}
	}
	return t
}
