(() => {
    const cnv = document.querySelector('canvas');
    const ctx = cnv.getContext('2d');

    let centerX = 0;
    let centerY = 0;
    function init() {
        cnv.width = innerWidth;
        cnv.height = innerHeight;

        centerX = cnv.width / 2;
        centerY = cnv.height / 2;
    }
    init();

    const numberOfRings = 17;
    const ringRadiusOffset = 15;
    const ringRadius = 150;
    const waveOffset = 10;
    const velocity = 1;
    let startAngles = Array.from({ length: numberOfRings }, () => Math.random() * 360);
    let ringVelocities = Array.from({ length: numberOfRings }, () => (Math.random() - 0.5) * 2 * velocity);
    let ringColors = Array.from({ length: numberOfRings }, () => {
        // Используем синие оттенки
        let h = Math.random() * 60 + 180; // оттенки от 180 до 240 градусов
        let s = Math.random() * 30 + 40; // насыщенность
        let l = Math.random() * 30 + 30; // светлота
        return `hsl(${h}, ${s}%, ${l}%)`;
    });

    function updateRings() {
        ctx.clearRect(0, 0, cnv.width, cnv.height); // Очистка экрана перед рисованием новых колец

        // Рисуем каждое кольцо
        for (let i = numberOfRings; i >= 1; i--) {
            let radius = i * ringRadiusOffset + ringRadius;
            let offsetAngle = i * waveOffset * Math.PI / 180;
            let alpha = i / numberOfRings;
            drawRing(radius, alpha, offsetAngle, i - 1);
        }
    }

    const maxWavesAmplitude = 10;
    const numberOfWaves = 10;

    function drawRing(radius, alpha, offsetAngle, ringIndex) {
        ctx.strokeStyle = 'rgba(255, 255, 255, 0)'; // Прозрачная обводка
        ctx.lineWidth = 9;

        let startAngle = startAngles[ringIndex];
        let ringVelocity = ringVelocities[ringIndex];
        let ringColor = ringColors[ringIndex];

        ctx.beginPath();

        for (let j = -180; j < 180; j++) {
            let currentAngle = (j + startAngle) * Math.PI / 180;
            let displacement = 0;
            let now = Math.abs(j);

            if (now > 30) {
                displacement = (now - 30) / 120;
            }

            if (displacement >= 1) {
                displacement = 1;
            }

            let angleOffset = Math.sin(startAngle * Math.PI / 180) * Math.PI / 180;
            let waveAmplitude = radius + displacement * Math.sin((currentAngle + angleOffset + offsetAngle) * numberOfWaves) * maxWavesAmplitude;
            let x = centerX + Math.cos(currentAngle) * waveAmplitude;
            let y = centerY + Math.sin(currentAngle) * waveAmplitude;
            j > -180 ? ctx.lineTo(x, y) : ctx.moveTo(x, y);
        }

        ctx.closePath();

        // Рисуем с использованием эффекта размытия
        ctx.save();
        ctx.globalCompositeOperation = 'lighter'; // Смешивание цветов
        ctx.fillStyle = ringColor; // Устанавливаем цвет заливки
        ctx.globalAlpha = 0.1; // Устанавливаем альфа-канал для смешивания
        ctx.fill(); // Заполняем цветом
        ctx.restore();

        ctx.stroke(); // Обводим линией

        // Обновляем startAngle для следующего кадра
        startAngles[ringIndex] = (startAngles[ringIndex] + ringVelocity) % 360;
    }

    function loop() {
        updateRings();
        requestAnimationFrame(loop);
    }
    loop();

    window.addEventListener('resize', init);
})();
