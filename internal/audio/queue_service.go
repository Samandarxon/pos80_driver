package audio

import (
	"log"
	"sync"
	"time"
)

// 🎯 AUDIO TASK STRUCTURE
type AudioTask struct {
	QueueNumber string
	RoomNumber  string
	Timestamp   time.Time
	Priority    int // 1 - High, 2 - Medium, 3 - Low
}

// 🚀 AUDIO QUEUE SERVICE
type AudioQueueService struct {
	audioService *AudioService
	tasks        chan AudioTask
	workerCount  int
	wg           sync.WaitGroup
	isRunning    bool
	mu           sync.RWMutex
}

func NewAudioQueueService(audioService *AudioService, workerCount int) *AudioQueueService {
	return &AudioQueueService{
		audioService: audioService,
		tasks:        make(chan AudioTask, 100), // 100 ta task uchun buffer
		workerCount:  workerCount,
		isRunning:    false,
	}
}

// 🎯 QUEUE NI ISHGA TUSHIRISH
func (q *AudioQueueService) Start() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.isRunning {
		return
	}

	q.isRunning = true

	// Worker larni ishga tushirish
	for i := 0; i < q.workerCount; i++ {
		q.wg.Add(1)
		go q.worker(i + 1)
	}

	log.Printf("🚀 Audio Queue Service started with %d workers", q.workerCount)
}

// 🛑 QUEUE NI TO'XTATISH
func (q *AudioQueueService) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if !q.isRunning {
		return
	}

	q.isRunning = false
	close(q.tasks)
	q.wg.Wait()

	log.Println("🛑 Audio Queue Service stopped")
}

// 📥 TASK QO'SHISH
func (q *AudioQueueService) AddTask(queueNumber, roomNumber string) {
	task := AudioTask{
		QueueNumber: queueNumber,
		RoomNumber:  roomNumber,
		Timestamp:   time.Now(),
		Priority:    2, // Default priority
	}

	select {
	case q.tasks <- task:
		log.Printf("📥 Audio task qo'shildi: %s -> %s (navbat: %d)",
			queueNumber, roomNumber, len(q.tasks))
	default:
		log.Printf("❌ Navbat to'la! Task qo'shilmadi: %s", queueNumber)
	}
}

// 👷 WORKER FUNCTION
func (q *AudioQueueService) worker(id int) {
	defer q.wg.Done()

	log.Printf("👷 Worker %d ishga tushdi", id)

	for task := range q.tasks {
		if !q.isRunning {
			break
		}

		log.Printf("🎯 Worker %d task bajarayapti: %s -> %s (qolgan: %d)",
			id, task.QueueNumber, task.RoomNumber, len(q.tasks))

		// Audio ni ijro etish
		if err := q.audioService.PlayAnnouncement(task.QueueNumber, task.RoomNumber); err != nil {
			log.Printf("❌ Worker %d xato: %v", id, err)
		} else {
			log.Printf("✅ Worker %d task tugatti: %s", id, task.QueueNumber)
		}

		// Keyingi task dan oldin qisqa pauza
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("👷 Worker %d to'xtadi", id)
}

// 📊 QUEUE STATUS
func (q *AudioQueueService) GetStatus() map[string]interface{} {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return map[string]interface{}{
		"is_running":   q.isRunning,
		"queue_length": len(q.tasks),
		"worker_count": q.workerCount,
		"buffer_size":  cap(q.tasks),
	}
}

// 🎯 TASK LARNI TOZALASH (agar kerak bo'lsa)
func (q *AudioQueueService) ClearQueue() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	count := 0
	for {
		select {
		case <-q.tasks:
			count++
		default:
			log.Printf("🗑️ Navbat tozalandi: %d task o'chirildi", count)
			return count
		}
	}
}
