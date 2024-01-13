document.getElementById('fileInput').addEventListener('change', function (e) {
    document.getElementById('imageName').innerHTML =
        document.getElementById('fileInput').files[0].name;
});

Fancybox.bind('[data-fancybox]', {
    // Your custom options
});
