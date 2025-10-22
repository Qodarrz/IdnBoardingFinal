# GreenFlow - Aplikasi Tracking Karbon User Sehari Hari

## Deskripsi Aplikasi  

*GreenFlow* adalah aplikasi website untuk memantau, mengurangi, dan memberi penghargaan atas jejak karbon harian pengguna. Aplikasi ini membantu pengguna memahami dampak aktivitas sehari-hari terhadap lingkungan dan mendorong gaya hidup yang lebih ramah lingkungan.


---

## Fitur Utama
- *Tracking Jejak Karbon:* Memantau aktivitas harian yang berdampak pada emisi karbon.
- *Saran Pengurangan Karbon:* Memberikan rekomendasi untuk mengurangi jejak karbon.
- *Penghargaan Pengguna:* Memberikan insentif bagi pengguna yang berhasil mengurangi jejak karbon.

---
## Teknologi
- *Backend:* Go dengan framework [Fiber](https://gofiber.io/) untuk performa tinggi dan skalabilitas.
- *Frontend:* React (libarray) dengan framework [Next.js](https://nextjs.org/) untuk pengalaman pengguna yang cepat dan interaktif.
- *Arsitektur:* Microservices, memungkinkan fitur dijalankan secara terpisah untuk memudahkan pengembangan dan monitoring.
---

## Arsitektur
GreenFlow mengikuti arsitektur monolith, sehingga semua modul seperti tracking, missions, dan reward berjalan dalam satu aplikasi terpadu. Hal ini membuat:

- Integrasi antar modul lebih sederhana karena berada dalam satu basis kode.
- Pengelolaan deployment lebih mudah karena hanya ada satu aplikasi yang dijalankan.
- Konsistensi data lebih terjaga karena semua modul berbagi basis data yang sama.

---

## Tim Pengembang
Muhammad Farhan Abdullah – UI/UX Designer

Renjie Syarbaini Prasetya – FrontEnd Developer

Muhammad Heidar Arrizqie – Backend Developer

Justine – DevOps

Kami merupakan perwakilan dari SMK Negeri 4 Kota Bogor.

# GreenFlow Backend Flow

Berikut adalah alur kerja backend GreenFlow dari request pengguna hingga response:

## 1. User Request
Pengguna melakukan request ke aplikasi, misalnya:
- Login
- Melakukan tracking aktivitas
- Order

## 2. Routing & Controller
Backend menerima request melalui routing dan diteruskan ke controller yang sesuai.

*Tugas Controller:*
- Memproses request
- Memanggil service yang relevan
- Menyiapkan response untuk pengguna

## 3. Service Layer
Service layer mengatur logika bisnis aplikasi.

*Contoh modul:*
- *Tracking:* Menghitung aktivitas pengguna  
- *Missions:* Memeriksa misi yang sedang aktif atau sudah selesai  
- *Rewards:* Mengecek reward yang bisa diberikan

## 4. Database Interaction
Service berkomunikasi dengan database untuk membaca atau menyimpan data pengguna, misi, atau reward.

*Keterangan:*
- Menggunakan model/ORM untuk mempermudah query  
- Validasi data dilakukan sebelum menyimpan ke database

## 5. Response ke User
Hasil proses dikirim kembali ke controller, kemudian diteruskan ke pengguna sebagai response API.

## 6. Error Handling & Logging
- Setiap langkah dicatat untuk debugging  
- Error ditangani agar aplikasi tetap stabil


---

## Fitur Utama

- *Chatbot:* Membantu pengguna dengan interaksi otomatis dan panduan cepat.  
- *Tracking GPS:* Memantau lokasi pengguna atau aktivitas secara real-time.  
- *Responsive Design:* Tampilan aplikasi menyesuaikan berbagai perangkat (desktop, tablet, mobile).  
- *Misi & Reward:* Pengguna dapat menyelesaikan misi dan mendapatkan reward.  
- *Notifikasi:* Memberikan update penting secara real-time kepada pengguna.  

---


## Teknologi yang Digunakan

- *Frontend:* React + Next.js + shadcn/ui (untuk komponen UI)  
- *Backend:* Golang dengan framework Fiber  
- *Database:* PostgreSQL  
- *Backend-as-a-Service / Auth:* Supabase  
- *Deployment:* Vercel (frontend)  
- *Otentikasi:* JWT (JSON Web Token)  
- *Upload File:* Langsung ke VPS
