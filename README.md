# EcoTrack - Lacak dan Kurangi Jejak Karbon Anda

**Slogan:** *“Langkah Kecilmu, Dampak Besar Bagi Bumi.”*

## Deskripsi Aplikasi
EcoTrack adalah aplikasi website untuk memantau, mengurangi, dan memberi penghargaan atas jejak karbon harian pengguna. Aplikasi ini membantu pengguna memahami dampak aktivitas sehari-hari terhadap lingkungan dan mendorong gaya hidup yang lebih ramah lingkungan.

## Fitur Utama
- **Tracking Jejak Karbon:** Memantau aktivitas harian yang berdampak pada emisi karbon.
- **Saran Pengurangan Karbon:** Memberikan rekomendasi untuk mengurangi jejak karbon.
- **Penghargaan Pengguna:** Memberikan insentif bagi pengguna yang berhasil mengurangi jejak karbon.

## Teknologi
- **Backend:** Go dengan framework [Fiber](https://gofiber.io/) untuk performa tinggi dan skalabilitas.
- **Frontend:** React (libarray) dengan framework [Next.js](https://nextjs.org/) untuk pengalaman pengguna yang cepat dan interaktif.
- **Arsitektur:** Microservices, memungkinkan fitur dijalankan secara terpisah untuk memudahkan pengembangan dan monitoring.

## Arsitektur
EcoTrack menggunakan arsitektur monolith, sehingga semua modul seperti tracking, missions, dan reward berjalan dalam satu aplikasi terpadu.  
Manfaat:
- Integrasi antar modul lebih sederhana karena berada dalam satu basis kode.
- Pengelolaan deployment lebih mudah karena hanya ada satu aplikasi yang dijalankan.
- Konsistensi data lebih terjaga karena semua modul berbagi basis data yang sama.

## Tim Pengembang
- **Muhamad Fauziansyah** – UI/UX Designer  
- **Justine** – Frontend Developer  
- **Muhammad Heidar Arrizqie** – Backend Developer  

Kami merupakan perwakilan dari SMK Negeri 4 Kota Bogor.

## Backend Flow
Alur kerja backend EcoTrack dari request pengguna hingga response:

1. **User Request**  
   Pengguna melakukan request seperti login, tracking aktivitas, atau order.

2. **Routing & Controller**  
   Backend menerima request dan diteruskan ke controller yang sesuai.  
   **Tugas Controller:**
   - Memproses request
   - Memanggil service yang relevan
   - Menyiapkan response untuk pengguna

3. **Service Layer**  
   Service layer mengatur logika bisnis aplikasi.  
   Contoh modul:
   - *Tracking:* Menghitung aktivitas pengguna
   - *Missions:* Memeriksa misi yang sedang aktif atau sudah selesai
   - *Rewards:* Mengecek reward yang bisa diberikan

4. **Database Interaction**  
   Service berkomunikasi dengan database untuk membaca atau menyimpan data.  
   - Menggunakan model/ORM untuk mempermudah query
   - Validasi data dilakukan sebelum menyimpan ke database

5. **Response ke User**  
   Hasil proses dikirim kembali ke controller, kemudian diteruskan ke pengguna sebagai response API.

6. **Error Handling & Logging**  
   - Setiap langkah dicatat untuk debugging  
   - Error ditangani agar aplikasi tetap stabil

## Fitur Tambahan
- **Chatbot:** Membantu pengguna dengan interaksi otomatis dan panduan cepat.
- **Tracking GPS:** Memantau lokasi pengguna atau aktivitas secara real-time.
- **Responsive Design:** Tampilan aplikasi menyesuaikan berbagai perangkat.
- **Misi & Reward:** Pengguna dapat menyelesaikan misi dan mendapatkan reward.
- **Notifikasi:** Memberikan update penting secara real-time kepada pengguna.

## Teknologi yang Digunakan
- **Frontend:** React + Next.js + shadcn/ui
- **Backend:** Golang dengan framework Fiber
- **Database:** PostgreSQL
- **Backend-as-a-Service / Auth:** Supabase
- **Deployment:** Vercel (frontend)
- **Otentikasi:** JWT (JSON Web Token)
- **Upload File:** Langsung ke VPS
