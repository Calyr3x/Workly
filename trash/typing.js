document.addEventListener('DOMContentLoaded', (event) => {
    document.querySelectorAll('.typewriter').forEach((element) => {
        const text = element.getAttribute('data-text');
        element.textContent = '';
        let i = 0;
        const speed = 150;

        function typeWriter() {
            if (i < text.length) {
                element.textContent += text.charAt(i);
                i++;
                setTimeout(typeWriter, speed);
            }
        }

        typeWriter();
    });
});
