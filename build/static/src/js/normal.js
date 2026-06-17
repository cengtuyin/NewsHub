/*  NORMAL.JS
*   全局通用方法
*   Made By Rexxrt
*/
let UI_Header;
window.onload = () => {
    _initPage();
    UI_Header = document.querySelector("html body Header");
    _checkScroll();
    window.onscroll = _checkScroll;
};



function _initPage() {
    AddMaterialIconClasses(`Tool>Icon`);
    AddMaterialIconClasses(`OptionsCard>Options>hOption>Icon`);
    try {
        document.querySelector("html>body>Header>Icon").onclick = () => {
            window.location.href = '/';
        };
    } catch (error) { }
    try {
        document.querySelector("Bottom>Level1>Icon").onclick = () => {
            window.location.href = '/';
        };
    } catch (error) { }
    try {
        const list = document.querySelectorAll('ul>li:has(input), Card>*:has(input[type="text"]), textarea');
        list.forEach(item => {
            const originalClick = item.onclick;
            const inputs = item.querySelectorAll('input, textarea');
            inputs.forEach(input => {
                const type = input.type;
                if (type === 'checkbox') {
                    const originalCheckboxClick = input.onclick;
                    item.addEventListener('click', (e) => {
                        e.stopPropagation();
                        input.checked = !input.checked;
                        const cardButtons = item.parentNode.parentNode.parentNode.lastElementChild;
                        if (cardButtons) cardButtons.removeAttribute('hide');
                        originalClick?.call(item);
                        originalCheckboxClick?.call(input, e);
                        for (let iii = 0; iii < 400 / 50; iii++) {
                            setTimeout(() => {
                                scrollToViewCenter(cardButtons)
                            }, iii * 50);
                        }
                    });

                } else if (type === 'text' || input.tagName == "TEXTAREA") {
                    item.addEventListener('click', (e) => {
                        e.stopPropagation();
                        originalClick?.call(item);
                    });

                    input.addEventListener('input', () => {
                        if (item.parentNode.tagName == 'CARD') item.parentNode.lastElementChild.removeAttribute('hide');
                        else item.parentNode.parentNode.parentNode.lastElementChild.removeAttribute('hide');
                    });
                }
            });
        });
    } catch (error) { }
    try {
        const list = document.querySelectorAll('OptionsCard');
        list.forEach(item => {
            const options = item.querySelectorAll('OptionsCard>Options>hOption');
            const views = item.querySelectorAll('View');
            options.forEach(option => {
                const view = views[option.style.getPropertyValue('--i')]
                const oevent = option.onclick;
                option.onclick = () => {
                    options.forEach(ooption => {
                        if (ooption == option) return;
                        ooption.removeAttribute('active');
                    });
                    views.forEach(oview => {
                        if (oview == view) return;
                        oview.setAttribute('hide', '');
                    });

                    option.setAttribute('active', '');
                    view.removeAttribute('hide');

                    if (oevent) oevent();
                };
            });
        });
    } catch (error) { }
    try {
        const list = document.querySelectorAll('SettingSetect');
        list.forEach(item => {
            const options_true = item.querySelectorAll('SettingSetect>Options')[0];
            item.addEventListener('wheel', function (e) {
                e.preventDefault();
                this.scrollLeft += e.deltaY;
            }, { passive: false });
            options_true.addEventListener('click', function (e) {
                if (e.target && e.target.matches('SettingSetect>Options>hOption,SettingSetect>Options>hOption>p,SettingSetect>Options>hOption>Icon,SettingSetect>Options>hOption>img')) {
                    let option = !e.target.matches('SettingSetect>Options>hOption') ? e.target.parentNode : e.target;
                    item.setect = option.getAttribute("tag");
                    if (option.getAttribute("tag") == "+" && item.getAttribute("add2left") != "") {
                        option = item.getAttribute("add2left") == "true" ? option.parentNode.firstElementChild : option.parentNode.lastElementChild;
                    }
                    options_true.querySelectorAll('SettingSetect>Options>hOption').forEach(ooption => {
                        if (ooption == option) return;
                        ooption.setAttribute('hide', '');
                    });
                    option.removeAttribute('hide');

                    item.setect = option.getAttribute("tag");
                    eval(`${item.getAttribute("onsetect").replace(/this/gi, "item")}`);
                }
            });
        });
    } catch (error) { }
}


function _checkScroll() {
    if (window.pageYOffset > 0) {
        UI_Header.classList.add("Header_Sub");
    } else {
        UI_Header.classList.remove("Header_Sub");
    }
}

function scrollToViewCenter(el) {
    const { top, height } = el.getBoundingClientRect();
    const elCenter = top + height / 2;
    const center = window.innerHeight / 2;
    window.scrollTo({
        top: (document.documentElement.scrollTop || document.body.scrollTop) - (center - elCenter),
        behavior: 'smooth'
    });
}

function AddMaterialIconClasses(selector) {
    const elements = document.querySelectorAll(selector);
    elements.forEach(el => {
        el.classList.add('material-icons', 'Icon');
    });
}

function ButtonLoadingAnimation(button, status) {
    if (status) button.setAttribute('loading', '');
    else button.removeAttribute('loading');
}

function CardButtonsDisplayMessage(button, status, message = "") {
    if (!status) button.parentNode.children[0].setAttribute('hide', '');
    else {
        button.parentNode.children[0].innerText = message;
        button.parentNode.children[0].removeAttribute('hide');
    }
}