package shm

/*
#include <stdlib.h>
#include <sys/shm.h>
#include <errno.h>
static int _errno() { return errno; }

typedef struct {
  key_t key;
  int shmid;
  char* data;
} _SHM;

static int create_shm(char* pathname, int size, _SHM *shm) {
  if ((shm->key = ftok(pathname, 'R')) == -1) {
    return 1;
  }
  if ((shm->shmid = shmget(shm->key, size, IPC_CREAT | 0666)) < 0) {
    return 1;
  }
  return 0;
}

static int attach_shm(char* pathname, int size, _SHM *shm) {
  if ((shm->key = ftok(pathname, 'R')) == -1) {
    return 1;
  }
  if ((shm->shmid = shmget(shm->key, size, 0666)) < 0) {
    return 1;
  }
  shm->data = (char *)shmat(shm->shmid, (void *)0, 0);
  if (shm->data == (char *)-1) {
    return 1;
  }
  return 0;
}

static int detatch_shm(_SHM *shm) {
  if (shmdt(shm->data) == -1){
    return 1;
  }
  return 0;
}

static int destroy_shm(_SHM *shm) {
  if (shmctl(shm->shmid, IPC_RMID, 0) == -1){
    return 1;
  }
  return 0;
}
*/
import "C"
import "errors"
import "syscall"
import "unsafe"

type Shm struct {
	shm C._SHM
}

func New(pathname string, size int) (*Shm, error) {
	p := C.CString(pathname)
	defer C.free(unsafe.Pointer(p))
	shm := &Shm{}
	if C.create_shm(p, C.int(size), &shm.shm) != 0 {
		return nil, errors.New(syscall.Errno(C._errno()).Error())
	}
	if C.attach_shm(p, C.int(size), &shm.shm) != 0 {
		return nil, errors.New(syscall.Errno(C._errno()).Error())
	}
	return shm, nil
}

func Attach(pathname string, size int) (*Shm, error) {
	p := C.CString(pathname)
	defer C.free(unsafe.Pointer(p))
	shm := &Shm{}
	if C.attach_shm(p, C.int(size), &shm.shm) != 0 {
		return nil, errors.New(syscall.Errno(C._errno()).Error())
	}
	return shm, nil
}

func (shm *Shm) Data() []byte {
	return (*[1 << 30]byte)(unsafe.Pointer(shm.shm.data))[:]
}

func (shm *Shm) Detatch() error {
	if C.detatch_shm(&shm.shm) != 0 {
		return errors.New(syscall.Errno(C._errno()).Error())
	}
	shm.shm.data = nil
	return nil
}

func (shm *Shm) Rm() error {
	if C.destroy_shm(&shm.shm) != 0 {
		return errors.New(syscall.Errno(C._errno()).Error())
	}
	shm.shm.key = -1
	shm.shm.shmid = -1
	shm.shm.data = nil
	return nil
}
