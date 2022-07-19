package script

func JSRegulation(abbreviation string) string {
	return `
	let regulation = {};
	regulation.regulation_name = document.getElementsByTagName('h1')[0].innerText;
	regulation.abbreviation = "` + abbreviation + `"
	JSON.stringify(regulation);
 `
}

var JSCheckChapter = `
	let h1s = document.getElementsByTagName("h1");
	if (h1s.length == 0) {
		false;
	} else {
		true;
	}
`

func JSChapter(regulationID, chapterOrderNum string) string {
	return `
	const terms = ["X","V", "I"];
	let chapter = {};
	chapter.regulation_id = ` + regulationID + `;
	chapter.order_num = ` + chapterOrderNum + `;
	chapter.chapter_name = "";
	chapter.chapter_num = "";
	// Chapter name
	let h1s = document.getElementsByTagName("h1");
    
	let h1 = "";
    
    if (h1s.length > 0){
        h1 = h1s[0].innerText;
        let aIdTags = h1s[0].getElementsByTagName("a");
        if (aIdTags && (aIdTags.length > 0)){
            let pId = aIdTags[0].id.split("dst")[1];
            chapter.chapter_id = parseInt(pId);
        }
    } 
    
	
	// Chapter num
	if (terms.some(term => h1.includes(term))){
		let h1_splited = h1.split(". ");
		chapter.chapter_num = h1_splited[0];
		chapter.chapter_name = h1_splited[1];
	} else {
		chapter.chapter_num = "";
		chapter.chapter_name = h1;
	}
	chapter.chapter_name = chapter.chapter_name.trim().replace((/  |\r\n|\n|\r/gm)," ").replace(/  +/g, ' ');
	h1s[0].remove();
	JSON.stringify(chapter);
	`
}

var JSCheckParagraphs = `
	let content = document.getElementsByClassName("document-page__content")[0];
	if (content){
		true;
	} else {
		false;
	}
`

func JSParagraphs(chapterID string) string {
	return `
	let chapter_id = ` + chapterID + `;
	let chapter = {};
	let paragraphs = [];
	let content = document.getElementsByClassName("document-page__content")[0];
	let i = 1;
	content.childNodes.forEach(function(el){
		let paragraph = {};
		paragraph.paragraph_order_num = i;
        paragraph.is_html = false;
        paragraph.is_table = false;
        paragraph.is_nft = false;
		paragraph.paragraph_class = "";
		paragraph.paragraph_text= "";
		paragraph.chapter_id = chapter_id;
        // If NTF
        if (el.attributes){
            for(let i = el.attributes.length - 1; i >= 0; i--) {
                if(el.attributes[i].name == "data-format-type" && el.attributes[i].value == "НФТ"){
                    // ID
                    let pIdTags = el.getElementsByTagName("a");
    				if (pIdTags && (pIdTags.length > 0)){
    					let pId = pIdTags[0].id.split("dst")[1];
    					paragraph.paragraph_id = parseInt(pId);
    					pIdTags[0].remove();
    				}
                    paragraph.is_nft = true;
                    paragraph.paragraph_text = el.innerHTML;
                    paragraphs.push(paragraph);
            		i++;
            		return;
                 }
            }
        }
       
        // If indent
		if (el.classList && el.classList.length == 1 && el.classList.contains("no-indent") && el.innerText.length == 0 ){
			paragraph.paragraph_text = "-";
            paragraph.paragraph_class = "indent";
            paragraphs.push(paragraph);
			i++;
			return;
		};
		
		// If table
		if (el.classList && el.classList.contains("doc-table")){
            paragraph.is_table = true;
			paragraph.paragraph_text = el.innerHTML;
            paragraphs.push(paragraph);
			i++;
			return;
		};
	
		// Drop useless elements
		if(el.textContent.length == 0){return}
		if (el.tagName != 'DIV' && el.tagName != 'P'){return;}
		
		if(el.classList && el.classList.contains("document__style")){
			let ps = el.getElementsByTagName("p");
			let str = "";
			Array.prototype.slice.apply(ps).forEach(function(e){
				let pIdTags = e.getElementsByTagName("a");
				if (pIdTags && (pIdTags.length > 0)){
					let pId = pIdTags[0].id.split("dst")[1];
					paragraph.paragraph_id = parseInt(pId);
					pIdTags[0].remove();
				}
				if (e.classList){
					paragraph.paragraph_class = e.classList[0];
				}
				str += e.innerText + " ";
			});
			if (str.length == 0){return;}
			paragraph.paragraph_text = str;
			paragraph.is_html = /<(?=.*? .*?\/ ?>|br|hr|input|!--|wbr)[a-z]+.*?>|<([a-z]+).*?<\/\1>/i.test(paragraph.paragraph_text); 
			paragraphs.push(paragraph);
			i++;
			return;
		};
		
		
		// Set class name
		if(el.classList && el.classList.contains("align_right")){
			// Paragraph Id
			let pIdTags = el.getElementsByTagName("a");
			if (pIdTags && (pIdTags.length > 0)){
				let pId = pIdTags[0].id.split("dst")[1];
				paragraph.paragraph_id = parseInt(pId);
				pIdTags[0].remove();
			}
	
			// drop info-link
			let divTag = el.getElementsByClassName("info-link")[0];
			if (divTag){divTag.remove()}
			console.log(el)
			paragraph.paragraph_class = "align_right";
            if (el.innerHTML.length == 0){return;}
			paragraph.paragraph_text = el.innerHTML;
			paragraph.is_html = /<(?=.*? .*?\/ ?>|br|hr|input|!--|wbr)[a-z]+.*?>|<([a-z]+).*?<\/\1>/i.test(paragraph.paragraph_text); 
			paragraphs.push(paragraph);
			i++;
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
            if (el.innerHTML.length == 0){return;}
			paragraph.paragraph_text = el.innerHTML;
			paragraph.is_html = /<(?=.*? .*?\/ ?>|br|hr|input|!--|wbr)[a-z]+.*?>|<([a-z]+).*?<\/\1>/i.test(paragraph.paragraph_text); 
			paragraphs.push(paragraph);
			i++;
			return;
		};
	
		// Set class name
		if(el.classList && el.classList.contains("align_left")){
			// Paragraph Id
			let pIdTags = el.getElementsByTagName("a");
			if (pIdTags && (pIdTags.length > 0)){
				let pId = pIdTags[0].id.split("dst")[1];
				paragraph.paragraph_id = parseInt(pId);
				pIdTags[0].remove();
			}
			
			// drop info-link
			let divTag = el.getElementsByClassName("info-link")[0];
			if (divTag){divTag.remove()}
		
			paragraph.paragraph_class = "align_left";
            if (el.innerHTML.length == 0){return;}
			paragraph.paragraph_text = el.innerHTML;
			paragraph.is_html = /<(?=.*? .*?\/ ?>|br|hr|input|!--|wbr)[a-z]+.*?>|<([a-z]+).*?<\/\1>/i.test(paragraph.paragraph_text); 
			paragraphs.push(paragraph);
			i++;
			return;
		};
	
		
		// Paragraph Id
		let pIdTags = el.getElementsByTagName("a");
        if (pIdTags && (pIdTags.length > 0)){
				let pId = pIdTags[0].id.split("dst")[1];
				paragraph.paragraph_id = parseInt(pId);
				pIdTags[0].remove();
		}
		
		let divTag = el.getElementsByClassName("info-link")[0];
		if (divTag){divTag.remove()}
		// Paragraph text
		// if inner contains links prepair them
		let links = Array.prototype.slice.apply( el.getElementsByTagName("a") );
		if (links.length > 0){
			links.forEach((el) => el.href = el.href.split("#dst")[1]);
		}
        if (el.innerHTML.length == 0){return;}
		paragraph.paragraph_text = el.innerHTML;
		i++;
		paragraph.is_html = /<(?=.*? .*?\/ ?>|br|hr|input|!--|wbr)[a-z]+.*?>|<([a-z]+).*?<\/\1>/i.test(paragraph.paragraph_text); 
		paragraphs.push(paragraph);
	});
	chapter.paragraphs = paragraphs;
	JSON.stringify(chapter).replace(/\\"/g, "\'");
`

}
