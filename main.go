package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	_ "github.com/golang/glog"
	_ "github.com/prometheus/client_golang/api"
	_ "github.com/prometheus/client_golang/api/prometheus/v1"
	_ "github.com/prometheus/common/model"
	"sync"
)

var (
	uuid string
)

func init() {
	flag.StringVar(&uuid, "uuid", "", "uuid")
	flag.Parse()
}

type HubService struct {
	mu sync.Mutex
}

func (s *HubService) Init() {

}

func main() {
	//k8s
	//k8s.PodsTest()
	//hubtest
	//ctx := gousb.NewContext()
	//defer ctx.Close()
	//hub := usb.HubController{
	//	Uuid:   "89LX0B5Z9",
	//	OperationType: "plugDevice",//
	//	UsbCtx: ctx,
	//}
	//hub.HandleControlDevice()
	//fmt.Println(twoSum([]int{3, 2, 4}, 6))
	fmt.Println(isPalindrome(121))

}

func twoSum(nums []int, target int) []int {
	valueToIndex := make(map[int]int)
	ret := make([]int, 0)
	for index, value := range nums {
		if secondValueIndex, exist := valueToIndex[target-value]; exist {
			if index < secondValueIndex {
				ret = append(ret, index)
				ret = append(ret, secondValueIndex)
			} else {
				ret = append(ret, secondValueIndex)
				ret = append(ret, index)
			}
			break
		}
		valueToIndex[value] = index
	}
	return ret
}


func isPalindrome(x int) bool {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, x)
	length := buf.Len()
	bytes := buf.Bytes()
	for index, value := range bytes {
		if bytes[length-1-index] != value {
			return false
		}
	}
	return true
}


 type ListNode struct {
	     Val int
	     Next *ListNode
 }

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1.Val == 0 && l1.Next == nil {
		return l2
	}
	if l2.Val == 0 && l2.Next == nil {
		return l1
	}

	result := &ListNode{}
	//迭代处理的结果节点，从第一个开始
	node := result
	// 进位
	var carry int
	// 和
	v := 0
	for {
		// 若 l1 还有，则累加，同时切换下一个l1
		if l1 != nil {
			v += l1.Val
			l1 = l1.Next
		}
		// 若 l2 还有，则累加，同时切换下一个 l2
		if l2 != nil {
			v += l2.Val
			l2 = l2.Next
		}
		//累加进位
		v += carry
		carry = v / 10 // 进位
		node.Val = v % 10 // 当前位
		// 若 l1 l2 进位都没了，结束
		if l1 == nil && l2 == nil && carry == 0 {
			break
		}
		// 构建下一个结果节点
		node.Next = &ListNode{}
		node = node.Next
		v = 0
	}

	return result
}