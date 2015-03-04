//有序列表
package tools

import (
	"container/list"
)

type SortedList struct {
	*list.List
	Limit       int                                       //最大列表个数
	compareFunc func(old, new interface{}) bool           //比较函数
	sameFunc    func(old, new, newValue interface{}) bool //相等比较函数
}

//建空链表
func CreateSortedList(limit int, compare func(old, new interface{}) bool, same func(old, new, newValue interface{}) bool) *SortedList {
	return &SortedList{list.New(), limit, compare, same}
}

func (this SortedList) findInsertPlaceElement(value interface{}) *list.Element {
	for element := this.Front(); element != nil; element = element.Next() {
		tempValue := element.Value
		if this.compareFunc(tempValue, value) {
			return element
		}
	}
	return nil
}

func (this SortedList) findNewPlaceElement(value interface{}, startPlace *list.Element) ([]interface{}, *list.Element) {
	var temp []interface{}
	var element *list.Element
	for element = startPlace.Prev(); element != nil; element = element.Prev() {
		tempValue := element.Value
		if this.compareFunc(tempValue, value) {
			temp = append(temp, tempValue)
		} else {
			return temp, element
		}
	}
	return temp, element
}

/**********************************************************
*排行榜中的某个值 数据发生变化 重新排行 返回超过的元素数组
*keyValue:标示这个值的key 用于sameFunc比较函数
*newValue:参与排行的新数据
***********************************************************/
func (this SortedList) ResortValue(keyValue, newValue interface{}) []interface{} {
	var temp []interface{}
	for element := this.Front(); element != nil; element = element.Next() {
		tempValue := element.Value
		if this.sameFunc(tempValue, keyValue, newValue) {
			var newplace *list.Element
			temp, newplace = this.findNewPlaceElement(tempValue, element)
			if len(temp) != 0 {
				if newplace == nil {
					this.List.PushFront(tempValue)
				} else {
					this.InsertAfter(tempValue, newplace)
				}
				this.List.Remove(element)
			}
			return temp
		}
	}
	return temp
}

/**********************************************************
*插入新值并排序(注意只能是新值)
***********************************************************/
func (this *SortedList) Insert(value interface{}) {
	if this.List.Len() == 0 {
		this.PushFront(value)
		return
	}
	if this.compareFunc(value, this.Back().Value) {
		if this.Len() < this.Limit {
			this.PushBack(value)
		}
		return
	}
	if this.compareFunc(this.List.Front().Value, value) {
		this.PushFront(value)
	} else {
		element := this.findInsertPlaceElement(value)
		if element != nil {
			this.InsertBefore(value, element)
		}
	}
	if this.Len() > this.Limit {
		this.Remove(this.Back())
	}
}
