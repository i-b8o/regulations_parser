js = `
// Chapter name
let chapterName = document.getElementsByTagName("h1")[0].innerText;
// Chapter num
let chapterNum = chapterName.split(". ")[0];
// Delete header, all indents  and info links
document.getElementsByTagName("h1")[0].remove();

Array.prototype.slice.apply( document.getElementsByClassName("info-link") ).map((el) => el.remove());
Array.prototype.slice.apply( document.getElementsByClassName("no-indent") ).map((el) => el.remove());
let paragraphs = [];

let pArray = Array.prototype.slice.apply( document.getElementsByTagName("p") );
pArray.map(function(p){
    let paragraph = {};
    // Paragraph Id
    let pIdTag = p.getElementsByTagName("a")[0];
    let pId = pIdTag.id.split("dst")[1];
    paragraph.id = parseInt(pId);
    pIdTag.remove();
    // Paragraph text
    // if inner contains links prepair them
    let links = Array.prototype.slice.apply( p.getElementsByTagName("a") );
    if (links.length > 0){
        links.forEach((el) => el.href = el.href.split("#dst")[1]);
    }

    paragraph.text = p.innerHTML;
    
    console.log(paragraph);
});

`