use std::env;

fn main() {
    let args: Vec<String> = env::args().collect();

    println!("open file {:?}", &args[1]);

    unsafe {
        let fd = libc::open(args[1].as_ptr() as *const libc::c_char, libc::O_RDWR, 0600);
        if fd == -1 {
            panic!("open failed: {:?}", std::io::Error::last_os_error());
        }

        let addr = libc::mmap(
            std::ptr::null_mut(),
            std::mem::size_of::<i32>(),
            libc::PROT_READ | libc::PROT_WRITE,
            libc::MAP_SHARED,
            fd,
            0,
        );
        if addr == libc::MAP_FAILED {
            panic!("mmap failed: {:?}", std::io::Error::last_os_error());
        }

        let j = addr as *mut *mut i32;
        let i = std::sync::atomic::AtomicPtr::from_ptr(j);

        loop {
            println!("{:?}", i.load(std::sync::atomic::Ordering::SeqCst) as i32);
            std::thread::sleep(std::time::Duration::from_micros(100));   
        }
    }
}
