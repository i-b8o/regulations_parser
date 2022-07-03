js = `
const terms = ["X","V", "I"];
let chapter = {};
// Chapter name
let h1s = document.getElementsByTagName("h1");
if (h1s.length > 0){
    chapter.name = h1s[0].innerText;
} else {
    chapter.name = "";
}
// Chapter num
if (terms.some(term => chapter.name.includes(term))){
    chapter.num = chapter.name.split(". ")[0];
} else {
    chapter.num = "";
}
// Delete header, all indents  and info links
document.getElementsByTagName("h1")[0].remove();
// Array.prototype.slice.apply( document.getElementsByClassName("info-link") ).map((el) => el.remove());
// Array.prototype.slice.apply( document.getElementsByClassName("no-indent") ).map((el) => el.remove());
let paragraphs = [];
let content = document.getElementsByClassName("document-page__content")[0];
content.childNodes.forEach(function(el){
    // If table
    if (el.classList && el.classList.contains("doc-table")){
        let paragraph = {};
        el.getElementsByTagName("table")[0].setAttribute('style', '');
        paragraph.text = el.innerHTML;
        paragraphs.push(paragraph);
        return;
    };
    
    // Drop useless elements
    if(el.textContent.length == 0){return}
    if (el.tagName != 'DIV' && el.tagName != 'P'){
            return;
    }
    
    if(el.classList && el.classList.contains("document__style")){return};
    if(el.classList && el.classList.contains("no-indent")){return};
    if(el.classList && el.classList.contains("document__format")){return};
    let paragraph = {};
    
    // Set class name
    if(el.classList && el.classList.contains("align_right")){
        paragraph.class = "align_right";
        paragraph.text = el.innerHTML;
        paragraphs.push(paragraph);
        return;
    };
    if(el.classList && el.classList.contains("align_center")){
        paragraph.class = "align_center";
        paragraph.text = el.innerHTML;
        paragraphs.push(paragraph);
        return;
    };
  
    // Paragraph Id
    let pIdTag = el.getElementsByTagName("a")[0];
    
    let pId = pIdTag.id.split("dst")[1];
    paragraph.id = parseInt(pId);
    pIdTag.remove();
    let divTag = el.getElementsByClassName("info-link")[0];
    if (divTag){divTag.remove()}
    // Paragraph text
    // if inner contains links prepair them
    let links = Array.prototype.slice.apply( el.getElementsByTagName("a") );
    if (links.length > 0){
        links.forEach((el) => el.href = el.href.split("#dst")[1]);
    }
    paragraph.text = el.innerHTML;
    
    paragraphs.push(paragraph);
});
chapter.paragraphs = paragraphs;
`
