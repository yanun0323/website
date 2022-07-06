# Strategy Pattern (策略模式)
> 定義演算法家族，個別封裝起來，讓他們之間可以互相替換。
> 此模式讓演算法的變動，不會影響到使用演算法的程式。

> 把會變動的部分取出並『封裝』起來，不讓其他部分部會受到影響。

</br>

在一個 `超類別` 當中，把會變動的 `行為` 可以另外再以 `interface` 封裝起來

這樣無論行為如何新增改變，都不會改動到 `超類別` 本身

</br>

## *範例：*

人類有`賺錢`及`花錢`兩種行為
```go
type Human struct {}
func(h *Human)MakeMoney(){ ... }
func(h *Human)SpendMoney(){ ... }
```
</br>

小明和小美都是人類，都會 `賺錢` `花錢`
```go
type Min struct {
    Human   /* 繼承人類型別 */
}
Min.MakeMoney()
Min.SpendMoney()

type Mei struct {
    Human   /* 繼承人類型別 */
}
Mei.MakeMoney()
Mei.SpendMoney()
```
</br>

但小華也是人類，但他是家裡蹲，並不會 `賺錢`
```go
type Hua struct {
    Human   /* 繼承人類型別 */
}
Mei.MakeMoney() /* 小華並不會賺錢，卻有賺錢行為，不合理 */
Mei.SpendMoney()
```
</br>

這時賺錢花錢的行為就是會變動的部分，可以另外封裝起來
```go
type Human struct {
    MakeMoneyBehavior IMakeMoneyBehavior   /* 封裝賺錢行為 */
    SpendMoneyBehavior ISpendMoneyBehavior    /* 封裝花錢行為 */
}
```

```go
/* 賺錢行為介面 */
type IMakeMoneyBehavior interface {
    MakeMoney()
}

/* 賺錢行為實作 */
type DoNotMakeMoney IMakeMoneyBehavior
func (d *DoNotMakeMoney) MakeMoney() {}

type MakeMoneyWithBrain IMakeMoneyBehavior
func (m *MakeMoneyWithBrain) MakeMoney() { ... }
```

```go
/* 花錢行為介面 */
type ISpendMoneyBehavior interface {
    SpendMoney()
}

/* 花錢行為實作 */
type DoNotSpentMoney ISpendMoneyBehavior
func (d *DoNotSpentMoney) SpendMoney() {}

type SpentMoneyForLover ISpendMoneyBehavior
func (s *SpentMoneyForLover) SpendMoney() { ... }
```