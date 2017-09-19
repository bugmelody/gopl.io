// http://127.0.0.1:6060/doc/codewalk/functions/

// http://cs.gettysburg.edu/projects/pig/piggame.html

/**
Game overview


(6-sided die 六面骰子)
Pig is a two-player game played with a 6-sided die. Each turn, you may roll or stay.

    If you roll a 1, you lose all points for your turn and play passes to your opponent.
	Any other roll adds its value to your turn score.

    If you stay, your turn score is added to your total score, and play passes to your opponent.

The first person to reach 100 total points wins.

The score type stores the scores of the current and opposing players,
in addition to the points accumulated during the current turn.
*/

package main

import (
	"fmt"
	"math/rand"
)

const (
	win            = 100 // The winning score in a game of Pig
	gamesPerSeries = 10  // The number of games per series to simulate
)

// A score includes scores accumulated in previous turns for each player,
// as well as the points scored by the current player in this turn.
type score struct {
	// 分别代表 player 的以前总得分, opponent 的以前总得分, player 本轮的得分
	player, opponent, thisTurn int
}

/**
User-defined function types
In Go, functions can be passed around just like any other value.(可以像值一样的被传递)
A function's type signature describes the types of its arguments and return values.
The action type is a function that takes a score and returns the resulting score and whether the current turn is over.

If the turn is over, the player and opponent fields in the resulting score
should be swapped,as it is now the other player's turn.
*/

// An action transitions(n. 过渡；转变) stochastically(adv. 偶估地；推测地；随机地) to a resulting score.
type action func(current score) (result score, turnIsOver bool)

/**
Multiple return values
Go functions can return multiple values.
The functions roll and stay each return a pair of values.
函数 roll 和 stay 均是返回两个值的函数
They also match the action type signature.
These action functions define the rules of Pig.
*/

// roll returns the (result, turnIsOver) outcome of simulating a die roll. （Die Roll: 掷骰）
// If the roll value is 1, then thisTurn score is abandoned, and the players'
// roles swap.  Otherwise, the roll value is added to thisTurn.
func roll(s score) (score, bool) {
	// outcome 代表了扔出来的骰子值
	outcome := rand.Intn(6) + 1 // A random int in [1, 6]
	if outcome == 1 {
		// player 扔出1,交换对手, player自己本轮得分清0
		return score{s.opponent, s.player, 0}, true
	}
	// player扔出2-6,仍然该自己, player 的 thisTurn 增加扔出筛子的数额
	return score{s.player, s.opponent, outcome + s.thisTurn}, false
}

// stay returns the (result, turnIsOver) outcome of staying.
// thisTurn score is added to the player's score, and the players' roles swap.
func stay(s score) (score, bool) {
	return score{s.opponent, s.player + s.thisTurn, 0}, true
}

/**
Higher-order functions
A function can use other functions as arguments and return values.
A strategy is a function that takes a score as input and returns an action to perform.
(Remember, an action is itself a function.)
*/
// A strategy chooses an action for any given score.
type strategy func(score) action

// stayAtK returns a strategy that rolls until thisTurn is at least k, then stays.
func stayAtK(k int) strategy {
	/**
	Function literals and closures
	Anonymous functions can be declared in Go, as in this example.
	Function literals are closures: they inherit the scope of the function in which they are declared.
	One basic strategy in Pig is to continue rolling until you have accumulated at least k points in a turn, and then stay.
	The argument k is enclosed by this function literal, which matches the strategy type signature.
	*/
	// 返回的匿名函数,它的函数签名匹配 strategy
	return func(s score) action {
		if s.thisTurn >= k {
			// 返回 stay 函数
			return stay
		}
		// 返回 roll 函数
		return roll
	}
}

/**
Simulating games
We simulate a game of Pig by calling an action to update the score until one player reaches 100 points.
Each action is selected by calling the strategy function associated with the current player.
*/
// play simulates a Pig game and returns the winner (0 or 1).
func play(strategy0, strategy1 strategy) int {
	// play 代表了两种策略之间的比赛

	// 声明并初始化一个 slice, 因为 strategies 长度为2,实际上 strategies[0] strategies[1] 代表了两个选手使用不同的策略
	strategies := []strategy{strategy0, strategy1}
	var s score
	var turnIsOver bool

	currentPlayer := rand.Intn(2) // [0,1) , 这里是要确定哪位选手轮先, 因为选手实际对应 strategy, 实际是随机哪个策略轮先
	// 后面也是一样, 实际上 currentPlayer==0,代表 strategy0, currentPlayer==1 ,代表 strategy1

	for s.player+s.thisTurn < win { // 循环条件: (之前分 + 本轮净得分) < 100, 每轮循环代表一个 turn

		// strategy: type strategy func(score) action
		action := strategies[currentPlayer](s) // action 要么是 stay, 要么是 roll
		s, turnIsOver = action(s)
		if turnIsOver {
			currentPlayer = (currentPlayer + 1) % 2 // 切换扔筛子权限
		}
	}


	// currentPlayer==0 实际代表 strategy0  ,currentPlayer==1 实际代表 strategy1
	return currentPlayer
}

/**
tally ['tælɪ] n. 计数器；标签；记账 vt. 使符合；计算；记录 vi. 一致；记分

Simulating a tournament
The roundRobin function simulates a tournament and tallies wins.
Each strategy plays each other strategy gamesPerSeries times.
*/
// roundRobin simulates a series of games between every pair of strategies.
func roundRobin(strategies []strategy) ([]int, int) {
	/**
	每两种策略之间 strategies[i], strategies[j], 会进行 gamesPerSeries 次比赛
	wins 是和 strategies 相同长度的slice
	wins[n] 代表了 strategies[n] 赢了多少次
	*/
	wins := make([]int, len(strategies))
	for i := 0; i < len(strategies); i++ {
		for j := i + 1; j < len(strategies); j++ {
			for k := 0; k < gamesPerSeries; k++ {
				winner := play(strategies[i], strategies[j])
				if winner == 0 {
					wins[i]++
				} else {
					wins[j]++
				}
			}
		}
	}

	gamesPerStrategy := gamesPerSeries * (len(strategies) - 1) // no self play
	// wins: 赢了多少次, gamesPerStrategy: 总共进行了多少次
	return wins, gamesPerStrategy
}

/**
Variadic function declarations(Variadic function:可变参数函数)
Variadic functions like ratioString take a variable number of arguments.
These arguments are available as a slice inside the function.
*/
// ratioString takes a list of integer values and returns a string that lists
// each value and its percentage of the sum of all values.
// e.g., ratios(1, 2, 3) = "1/6 (16.7%), 2/6 (33.3%), 3/6 (50.0%)"
func ratioString(vals ...int) string {
	// 开始计算 total
	total := 0
	for _, val := range vals {
		total += val
	}
	// total 计算完毕

	// 开始计算要返回的字符串
	s := ""
	for _, val := range vals {
		if s != "" {
			s += ", "
		}
		pct := 100 * float64(val) / float64(total)
		s += fmt.Sprintf("%d/%d (%0.1f%%)", val, total, pct)
	}
	return s
}

/**
Simulation(n. 仿真；模拟；模仿；假装) results
The main function defines 100 basic strategies, simulates a round robin
tournament, and then prints the win/loss record of each strategy.
Among these strategies, staying at 25 is best, but the optimal strategy
for Pig is much more complex.
*/
func main() {
	// 100 种策略
	strategies := make([]strategy, win)

	for k := range strategies {
		// 每种策略在得到k分后 stay, 这样最后我们统计到得到多少分就stay是最优选择
		strategies[k] = stayAtK(k + 1)
	}

	wins, games := roundRobin(strategies)

	for k := range strategies {

		/**
		wins[k] 代表 stayAtK 策略赢了多少次
		games 代表 stayAtK 策略总共玩了多少次
		games-wins[k] 代表输了多少次
		 */
		fmt.Printf("Wins, losses staying at k =% 4d: %s\n",
			k+1, ratioString(wins[k], games-wins[k]))
	}
}

/**
输出结果
          stayAtK(k)这种策略:     赢/总 (百分比),    输/总(百分比)
Wins, losses staying at k =   1: 241/990 (24.3%), 749/990 (75.7%)
Wins, losses staying at k =   2: 243/990 (24.5%), 747/990 (75.5%)
 */
