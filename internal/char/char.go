/*
Copyright © 2022 Aurelio Calegari, et al.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package char

type Char struct {
	Coordinates []Coordinates
	Next        int
	PaintChar   rune
	Shade       rune
}

type Coordinates struct {
	X int
	Y int
	L int
}

func (c *Char) GetWidth() int {
	maxWidth := 0
	for _, v := range c.Coordinates {
		if v.X+v.L > maxWidth {
			maxWidth = v.X + v.L
		}
	}
	// account shade
	maxWidth = maxWidth + 2
	if c.Next > maxWidth {
		return c.Next
	}
	return maxWidth
}

var LoggoLogo = []Char{CharacterL, CharacterApostrophe, CharacterO, CharacterRevG, CharacterG, CharacterO}

/*
 012345678
0╦╦╦╦╦╦╦╬
1▓▓▓░╬╬╬╬
2╬╬▓▓░╬╬╬
3╬╬▓▓░╬╬╬
4╬╬▓▓░╬╬╬
5╬╬▓▓░╬╬╬
6╬╬▓▓░╬╬╬
7╬╬▓▓░╬╬╬
8╬╬▓▓░╬╬╬
9╬╬╬▓▓▓░╬╳
*/

var CharacterL = Char{
	PaintChar: '▓',
	Shade:     '░',
	Coordinates: []Coordinates{
		{0, 1, 3},
		{2, 2, 2},
		{2, 3, 2},
		{2, 4, 2},
		{2, 5, 2},
		{2, 6, 2},
		{2, 7, 2},
		{2, 8, 2},
		{3, 9, 3},
	},
	Next: 8,
}

/*
 012
0▓▓░
1╬▓░
2╬╬╬
3╬╬╬
4╬╬╬
5╬╬╬
6╬╬╬
7╬╬╬
8╬╬╬
9╳╬╬
*/

var CharacterApostrophe = Char{
	PaintChar: '▓',
	Shade:     '░',
	Coordinates: []Coordinates{
		{0, 0, 2},
		{0, 1, 1},
	},
	Next: 0,
}

/*
 01234567890123
0╦╦╦╦╦╦╦╦╦╦╦╦
1╬╬╬╬╬╬╬╬╬╬╬╬
2╬╬╬╬╬╬╬╬╬╬╬╬
3╬╬╬╬╬╬╬╬╬╬╬╬
4╬╬╬▓▓▓▓▓░╬╬╬
5╬▓▓░╬╬╬╬▓▓░╬
6▓▓░╬╬╬╬╬╬▓▓░
7▓▓░╬╬╬╬╬╬▓▓░
8╬▓▓░╬╬╬╬▓▓░╬
9╬╬╬▓▓▓▓▓░╬╬╬╬╳
*/

var CharacterO = Char{
	PaintChar: '▓',
	Shade:     '░',
	Coordinates: []Coordinates{
		{3, 4, 5},

		{1, 5, 2},
		{0, 6, 2},
		{0, 7, 2},
		{1, 8, 2},

		{8, 5, 2},
		{9, 6, 2},
		{9, 7, 2},
		{8, 8, 2},

		{3, 9, 5},
	},
	Next: 13,
}

/*
 012345678901234567
0╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦
1╬╬╬╬╬▓▓▓▓▓▓░╬╬╬╬
2╬╬╬▓▓░╬╬╬╬▓▓▓░╬╬
3╬╬▓▓░╬╬╬╬╬╬╬╬╬╬╬
4╬▓▓░╬╬╬╬╬╬╬╬╬╬╬╬
5▓▓░╬╬╬╬╬╬╬╬╬╬╬╬╬
6▓▓░╬╬╬╬╬╬╬▓▓▓▓▓░
7▓▓░╬╬╬╬╬╬╬╬╬▓▓░╬
8╬╬▓▓░╬╬╬╬╬╬▓▓▓░╬
9╬╬╬╬▓▓▓▓▓▓▓░╬╬╬╬╬╳
*/

var CharacterG = Char{
	PaintChar: '▓',
	Shade:     '░',
	Coordinates: []Coordinates{
		{5, 1, 6},
		{3, 2, 2},
		{10, 2, 3},
		{2, 3, 2},
		{1, 4, 2},
		{0, 5, 2},
		{0, 6, 2},
		{0, 7, 2},
		{4, 9, 7},
		{10, 6, 5},
		{12, 7, 2},
		{11, 8, 3},
		{2, 8, 2},
	},
	Next: 17,
}

/*
 012345678901234567
0╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦╦
1╬╬╬╬▓▓▓▓▓▓░╬╬╬╬╬
2╬╬▓▓▓░╬╬╬╬▓▓░╬╬╬
3╬╬╬╬╬╬╬╬╬╬╬▓▓░╬╬
4╬╬╬╬╬╬╬╬╬╬╬╬▓▓░╬
5╬╬╬╬╬╬╬╬╬╬╬╬╬▓▓░
6▓▓▓▓▓░╬╬╬╬╬╬╬▓▓░
7╬▓▓░╬╬╬╬╬╬╬╬╬▓▓░
8╬▓▓▓░╬╬╬╬╬╬▓▓░╬╬
9╬╬╬╬▓▓▓▓▓▓▓░╬╬╬╬╬╳
*/

var CharacterRevG = Char{
	PaintChar: '▓',
	Shade:     '░',
	Coordinates: []Coordinates{
		{4, 1, 6},
		{2, 2, 3},
		{0, 6, 5},
		{1, 7, 2},
		{1, 8, 3},
		{4, 9, 7},
		{10, 2, 2},
		{11, 3, 2},
		{12, 4, 2},
		{13, 5, 2},
		{13, 6, 2},
		{13, 7, 2},
		{11, 8, 2},
	},
	Next: 17,
}
