jsRegulation = `
    let regulation = {};
    regulation.regulation_name = document.getElementsByTagName('h1')[0].innerText;
    JSON.stringify(regulation);
`

jsChapter = `
const terms = ["X","V", "I"];
let chapter = {};
// Chapter name
let h1s = document.getElementsByTagName("h1");
let h1 = h1s.length > 0 ? h1s[0].innerText : "";

// Chapter num
if (terms.some(term => h1.includes(term))){
    let h1_splited = h1.split(". ");
    chapter.chapter_num = h1_splited[0];
    chapter.chapter_name = h1_splited[1];
} else {
    chapter.chapter_num = "";
    chapter.chapter_name = h1;
}
// Delete header, all indents  and info links
h1s[0].remove();
`
jsParagraphs = `
let paragraphs = [];
let content = document.getElementsByClassName("document-page__content")[0];
content.childNodes.forEach(function(el){
    // If table
    if (el.classList && el.classList.contains("doc-table")){
        let paragraph = {};
        el.getElementsByTagName("table")[0].setAttribute('style', '');
        paragraph.paragraph_text = el.innerHTML;
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
         // Paragraph Id
        let pIdTags = el.getElementsByTagName("a");
        if (pIdTags && (pIdTags.length > 0)){
            let pId = pIdTags[0].id.split("dst")[1];
            paragraph.paragraph_id = parseInt(pId);
            pIdTags[0].remove();           
        }
        
        let divTag = el.getElementsByClassName("info-link")[0];
        if (divTag){divTag.remove()}

        paragraph.paragraph_class = "align_right";
        paragraph.paragraph_text = el.innerHTML;
        paragraphs.push(paragraph);
        return;
    };
    if(el.classList && el.classList.contains("align_center")){
        // Paragraph Id
        let pIdTags = el.getElementsByTagName("a");
        if (pIdTags && (pIdTags.length > 0)){
            let pId = pIdTags[0].id.split("dst")[1];
            paragraph.paragraph_id = parseInt(pId);
            pIdTags[0].remove();           
        }
        
        let divTag = el.getElementsByClassName("info-link")[0];
        if (divTag){divTag.remove()}
        
        paragraph.paragraph_class = "align_center";
        paragraph.paragraph_text = el.innerHTML;
        paragraphs.push(paragraph);
        return;
    };
  
    // Paragraph Id
    let pIdTag = el.getElementsByTagName("a")[0];
    
    let pId = pIdTag.id.split("dst")[1];
    paragraph.paragraph_id = parseInt(pId);
    pIdTag.remove();
    let divTag = el.getElementsByClassName("info-link")[0];
    if (divTag){divTag.remove()}
    // Paragraph text
    // if inner contains links prepair them
    let links = Array.prototype.slice.apply( el.getElementsByTagName("a") );
    if (links.length > 0){
        links.forEach((el) => el.href = el.href.split("#dst")[1]);
    }
    paragraph.paragraph_text = el.innerHTML;
    
    paragraphs.push(paragraph);
});
chapter.paragraphs = paragraphs;
JSON.stringify(chapter).replace(/\\"/g, "\'");
`
