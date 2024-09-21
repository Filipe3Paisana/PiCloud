function downloadFile() {
    const link = document.createElement('a');
    link.href = 'path/to/your/file'; // Substitua pelo caminho real do arquivo
    link.download = 'filename'; // Substitua pelo nome real do arquivo
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}
