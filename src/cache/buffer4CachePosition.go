package cache

import (
	//	"time"
	"errors"
	//	"time"
	"fmt"
)

// 缓存的数量
const BUFFER4CACHEPOSITION = 3

// 写入cache的缓冲区，用于一次性写入多少条记录，而非一条一条间隔写
type Buffer4CachePosition struct {
	idxWrite chan int          //	标识写哪个缓冲
	idxRead  chan int          // 标识读哪个缓冲
	buf1     chan UserPosition // 第一个缓冲区
	buf2     chan UserPosition // 第二个缓冲区
	num1     int               // 标识在第一个缓冲区中写了几个数据
	num2     int               // 标识在第二个缓冲区中写了几个数据
}

// 初始化
func InitBuffer4CachePosition() *Buffer4CachePosition {
	pRet := new(Buffer4CachePosition)

	pRet.idxWrite = make(chan int, 1)
	pRet.idxRead = make(chan int, 1)
	pRet.buf1 = make(chan UserPosition, BUFFER4CACHEPOSITION)
	pRet.buf2 = make(chan UserPosition, BUFFER4CACHEPOSITION)
	pRet.num1 = 0
	pRet.num2 = 0

	// 初始向第一个缓冲区中写入数据
	pRet.idxWrite <- 1

	return pRet
}

// 向缓冲中写入数据
// chanPosition－－cache接受数据的channel
func (pBuf *Buffer4CachePosition) WriteData2Buffer(pSet *RealTimePositionSet, chanPosition <-chan UserPosition) error {
	fmt.Println("---------------位置WriteData2Buffer：")

	for {
		idx := <-pBuf.idxWrite
		fmt.Println("---------------位置WriteData2Buffer-----------idx：", idx)
		if idx == 1 {
			for pos := range chanPosition {
				//fmt.Println("---------------------pos1", pos)
				pBuf.buf1 <- pos // = append(pBuf.buf1, pos)
				pBuf.num1++
				if pBuf.num1 >= BUFFER4CACHEPOSITION {
					pBuf.num1 = 0
					//pSet.idx <- 1
					pBuf.idxWrite <- 2
					fmt.Println("1---xxx")
					pBuf.idxRead <- 1
					fmt.Println("1---yyy")
					break
				}
			}
			/*for{
				select{
					case pos := <- chanPosition:
						pBuf.buf1 <- pos
						pBuf.num1++
						if pBuf.num1 >= BUFFER4CACHEPOSITION{
							pBuf.num1 = 0
							pBuf.idxWrite <- 2
							pBuf.idxRead <- 1
							break
						}
					case <- time.After(time.Second * 5):
						pBuf.num1 = 0
						pBuf.idxWrite <- 2
						pBuf.idxRead <- 1
						break
				}
			}*/
		} else {
			for pos := range chanPosition {
				//fmt.Println("---------------------pos2", pos)
				pBuf.buf2 <- pos // = append(pBuf.buf2, pos)
				pBuf.num2++
				if pBuf.num2 >= BUFFER4CACHEPOSITION {
					pBuf.num2 = 0
					//pSet.idx <- 1
					pBuf.idxWrite <- 1
					fmt.Println("2---xxx")
					pBuf.idxRead <- 2
					fmt.Println("2---yyy")
					break
				}
			}
			/*for{
				select{
					case pos := <- chanPosition:
						pBuf.buf2 <- pos
						pBuf.num2++
						if pBuf.num1 >= BUFFER4CACHEPOSITION{
							pBuf.num2 = 0
							pBuf.idxWrite <- 1
							pBuf.idxRead <- 2
							break
						}
					case <- time.After(time.Second * 5):
						pBuf.num2= 0
						pBuf.idxWrite <- 1
						pBuf.idxRead <- 2
						break
				}
			}*/
		}
	}
	fmt.Println("---------------位置WriteData2Buffer结束：")
	return errors.New("error in writing data to buffer for cache")
}

// 写channel数据至slice
func writeChannel2Slice(backslice []UserPosition, slice chan<- UserPosition, buf <-chan UserPosition) {
	i := 0
	for v := range buf {
		//fmt.Println("channel -- slice", v)
		slice <- v
		//slice = append(slice, v)
		backslice = append(backslice, v)

		i++
		if i >= BUFFER4CACHEPOSITION {
			break
		}
	}
	/*for{
		select{
			case v := <- buf:
				slice <- v
				backslice = append(backslice, v)
				i++
				if i >= BUFFER4CACHEPOSITION{
					break
				}
			case <-time.After(time.Second * 3):
				break

		}
	}*/
}

// 写入cache存储中，同时写入db的缓冲区
func (pBuf *Buffer4CachePosition) Write2CacheAndDbBuffer(pDbBuf *Buffer4Db, pRealSet *RealTimePositionSet) error {

	fmt.Println("--------------------Write2CacheAndDbBuffer")
	for {
		idx := <-pBuf.idxRead

		// 写入db缓存
		//go func(pDbBuf *Buffer4Db) {
		// 写入数据库缓存
		backslice := make([]UserPosition, 0)

		for {
			if idx == 1 {
				writeChannel2Slice(backslice, pDbBuf.chBuf, pBuf.buf1)
				break
			} else {
				writeChannel2Slice(backslice, pDbBuf.chBuf, pBuf.buf2)
				break
			}
			/*if pDbBuf.idx == 1 {
				//go func(pDbBuf *Buffer4Db) {
				fmt.Println("--------------------写入dbbuf1缓存")
				if idx == 1 {
					writeChannel2Slice(backslice, pDbBuf.buf1, pBuf.buf1)
					//pDbBuf.buf1 = append(pDbBuf.buf1, pBuf.buf1...)
				} else {
					writeChannel2Slice(backslice, pDbBuf.buf1, pBuf.buf2)
					//pDbBuf.buf1 = append(pDbBuf.buf1, pBuf.buf2...)
				}
				fmt.Println("pDbBuf.buf1 len: ", len(pDbBuf.buf1))
				if len(pDbBuf.buf1) >= BUFFERSIZE4DB {
					fmt.Println("3---xxx")
					pDbBuf.chIdx <- 1
					fmt.Println("3---yyy")
					pDbBuf.idx = 2
					fmt.Println("3---zzz")
					break
				}

			} else {
				fmt.Println("--------------------写入dbbuf2缓存")
				if idx == 1 {
					writeChannel2Slice(backslice, pDbBuf.buf2, pBuf.buf1)
				} else {
					writeChannel2Slice(backslice, pDbBuf.buf2, pBuf.buf2)
				}
				fmt.Println("pDbBuf.buf2 len: ", len(pDbBuf.buf1))
				if len(pDbBuf.buf2) >= BUFFERSIZE4DB {
					fmt.Println("4---xxx")
					pDbBuf.chIdx <- 2
					fmt.Println("4---yyy")
					pDbBuf.idx = 1
					fmt.Println("4---zzz")
					break
				}
			}*/
		}
		//}(pDbBuf)

		// 写入cache
		pRealSet.RwLock.Lock()

		tmpSlice := backslice
		/*tmpSlice := make([]UserPosition, 0)
		if idx == 1 {
			tmpSlice = append(tmpSlice, pBuf.buf1...)
			pBuf.buf1 = pBuf.buf1[:0]
		} else {
			tmpSlice = append(tmpSlice, pBuf.buf2...)
			pBuf.buf2 = pBuf.buf2[:0]
		}*/

		size := len(pRealSet.allPositions)
		//pRealSet.allPositions = append(pRealSet.allPositions, pBuf.buf1...)
		for i, elm := range tmpSlice {
			elm.Id = int64(len(pRealSet.allPositions))
			pRealSet.allPositions = append(pRealSet.allPositions, elm)

			// 添加用户id到轨迹的映射
			if _, exist := pRealSet.userId2Postions[elm.UserId]; !exist {
				pRealSet.userId2Postions[elm.UserId] = make(UserPositionIds, 0)
			}
			pRealSet.userId2Postions[elm.UserId] = append(pRealSet.userId2Postions[elm.UserId], int64(size+i))

			// 添加时间到位置的映射
			if _, exist := pRealSet.time2Positions[elm.TmSpace.CreateTime]; !exist {
				pRealSet.time2Positions[elm.TmSpace.CreateTime] = make(UserPositionIds, 0)
			}
			pRealSet.time2Positions[elm.TmSpace.CreateTime] = append(pRealSet.time2Positions[elm.TmSpace.CreateTime], int64(size+i))
		}
		pRealSet.RwLock.Unlock()
	}
	fmt.Println("--------------------Write2CacheAndDbBuffer结束")
	return errors.New("error in write buffer to cache and db buffer")
}
