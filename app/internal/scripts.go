package main

var jsRegulation = `
	let regulation = {};
	regulation.regulation_name = document.getElementsByTagName('h1')[0].innerText;
	JSON.stringify(regulation);
 `

func jsChapter(chapterID string) string {
	return `
	let chapter_id = ` + chapterID + `;
	let chapter = {};
	let paragraphs = [];
	let content = document.getElementsByClassName("document-page__content")[0];
	let i = 0;
	content.childNodes.forEach(function(el){
		let paragraph = {};
		paragraph.paragraph_id = 0;
		paragraph.paragraph_order_num = i;
		paragraph.paragraph_class = "";
		paragraph.paragraph_text= "";
		paragraph.chapter_id = chapter_id;
	
		
		// If table
		if (el.classList && el.classList.contains("doc-table")){
			el.getElementsByTagName("table")[0].setAttribute('style', '');
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
				str += e.innerText;
			});
			paragraph.paragraph_text = str;
			paragraphs.push(paragraph);
			i++;
			return;
		};
		//if(el.classList && el.classList.contains("no-indent")){return};
		//if(el.classList && el.classList.contains("document__format")){return};
		
		
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
			paragraph.paragraph_text = el.innerHTML;
			paragraphs.push(paragraph);
			i++;
			console.log(paragraph)
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
			paragraph.paragraph_text = el.innerHTML;
			paragraphs.push(paragraph);
			i++;
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
		i++;
		paragraphs.push(paragraph);
	});
	chapter.paragraphs = paragraphs;
	JSON.stringify(chapter).replace(/\\"/g, "\'");
	`

}
