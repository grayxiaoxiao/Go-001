学习笔记
[作业地址](https://github.com/Go-000/Go-000/issues/8)

##### 前言
> 对微服务的概念不是很清晰，对其架构分层也不清晰。通过查询资料大致知道DAO层的主要作用是沟通数据，并且映射数据到对应的代码结构上，以此作为服务调用的基础。
> 
> **理解有误指出还请指正。**

#### 作业内容
> 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

##### 个人看法
> `sql.ErrNoRows`是对`QueryRow`或者`QueryRowContext`的结果进行`Scan`操作时，query的结果没有任何数据的error信息。

**我觉着是否需要将这个error信息`wrap`给上层调用者，大约可以分为下面两种：**
###### 1. 在根据具有唯一性标识的条件获取单条数据时（例如查看指定文章详情的操作），产生的`sql.ErrNoRows`需要`wrap`给上层调用者。
> 这种情况下说明给与的唯一标识参数可能存在错误或者数据由于多线操作已被移除等等，为了避免空数据参与后续的业务操作，所以需要明确告知上层调用者该错误信息，同时也阻止由于`nil`之类的数值参与后续代码的计算导致其他错误。

```golang
// DAO
package model
type Order struct {
	Serial string `json:"serial"` // serial作为唯一标识
	// 其他字段...
}

fun GetOrderBySerial(serial string) order, error {
	// connect DB ...
	var order Order
	err := sql.QueryRowContext(ctx, "select serial, .. from orders where serial = ?", serial).Scan(&order.Serial,...)
	switch {
		case err == sql.ErrNoRows:
			error_message := "Not find order with serial = " + serial
			return nil, errors.Wrap(err, error_message)
		case err != nil:
			error_message := "Other error with query!!"
			return nil, errors.Wrap(err, error_message)
		default:
			return order, nil
	}
	
}

// Service
package service

func GetOrderBySerial(serial string) json {
	order, err := model.GetOrderBySerial(serial)
	if err == sql.ErrNoRows {
		render_json({code: -1, message: "NOT FOUND OR Illegal parameter", data: {}})
	}
	render_json({code: 0, message: "Success", data: order.to_json})
}
```

###### 2. 在获取某一类数据列表时（例如查看热点视频推荐的操作），产生的`sql.ErrNoRows`不需要`wrap`给上层调用者。
> 这种情况下说明数据中并没有指定类型的数据集。这种情况下，给上层调用者返回空数据集合类型即可（如空数组）。通常情况下，此类情况下后续的数据集操作并不具有很高的危险性。

```golang
// DAO
package model
type Order struct {
	Serial string `json:"serial"` // 编号，视作为唯一标识
	Amount int64  `json:"amount"` // 金额
	//其他字段..
}

func GetOrdersWithCondition(condition map[string]interface{}) []Order, error {
	// connect DB ...
	rows, err := sql.QueryRowContext(ctx, "select serial, amount, ... from orders where amount > ?", condition["amount"])
	if err != nil {
		return [], err
	}
	
	orders := make([]Order, 10)
	for rows.Next() {
		var order Order
		err = rows.Scan(&order.serial, &order.amount, ..)
		if err == sql.ErrNoRows {
			// break
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// Service
package service

func GetOrdersWithCondition(condition map[string]interface{}) json {
	orders, err := model.GetOrdersWithCondition(condition)
	if err != nil {
		log.Fatal(err)
		render_json({code: -2, message: err.Error(), data: []})
	}
	render_json({code: 0, message: "Success", data: orders.to_json})
}
```
