
const documents = [
    { name: "Documento1.pdf", id: 1 },
    { name: "Foto1.png", id: 2 },
    { name: "Relatorio.docx", id: 3 },
];

function loadFiles() {
    const fileList = document.getElementById('fileList');
    fileList.innerHTML = ''; 

    documents.forEach(doc => {
        const li = document.createElement('li');

        
        const docName = document.createElement('span');
        docName.textContent = doc.name;
        docName.style.marginRight = "20px";
        
        const downloadButton = document.createElement('button');
        downloadButton.textContent = "Download";
        downloadButton.onclick = () => confirmDownload(doc.name);

        
        li.appendChild(docName);
        li.appendChild(downloadButton);
        fileList.appendChild(li);
    });
}

function confirmDownload(fileName) {
    if (confirm(`Deseja fazer download de ${fileName}?`)) {
        downloadFile(fileName);
    }
}

function downloadFile(fileName) {
    alert(`Iniciando o download de: ${fileName}`);
    
}


loadFiles();


