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

const _initPage_Icon            = 0b1;
const _initPage_HanderIcon      = 0b10;
const _initPage_BottomIcon      = 0b100;
const _initPage_CardButtons     = 0b1000;
const _initPage_OptionsCard     = 0b10000;
const _initPage_SettingSetect   = 0b100000;

function _initPage(level = 0b1111111111111111111111111) {
    if ((level & _initPage_Icon) != 0) {
        AddMaterialIconClasses(`Tool>Icon`);
        AddMaterialIconClasses(`OptionsCard>Options>hOption>Icon`);
    }
    if ((level & _initPage_HanderIcon) != 0) {
        try {
            document.querySelector("html>body>Header>Icon").onclick = () => {
                window.location.href = '/';
            };
        } catch (error) { }
    }
    if ((level & _initPage_BottomIcon) != 0) {
        try {
            document.querySelector("Bottom>Level1>Icon").onclick = () => {
                window.location.href = '/';
            };
        } catch (error) { }
    }
    if ((level & _initPage_CardButtons) != 0) {
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
    }
    if ((level & _initPage_OptionsCard) != 0) {
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
    }
    if ((level & _initPage_SettingSetect) != 0) {
        try {
            const list = document.querySelectorAll('SettingSetect');
            list.forEach(item => {
                const options_true = item.querySelectorAll('SettingSetect>Options')[0];
                item.addEventListener('wheel', function (e) {
                    e.preventDefault();
                    this.scrollLeft += e.deltaY / 2;
                }, { passive: false });
                options_true.addEventListener('click', function (e) {
                    if (e.target && e.target.matches('SettingSetect>Options>hOption,SettingSetect>Options>hOption>p,SettingSetect>Options>hOption>Icon,SettingSetect>Options>hOption>img')) {
                        let option = !e.target.matches('SettingSetect>Options>hOption') ? e.target.parentNode : e.target;
                        item.setect = option.getAttribute("tag");
                        if (option.getAttribute("tag") == "+" && item.getAttribute("add2left") != "") {
                            option = item.getAttribute("add2left") == "true" ? option.parentNode.firstElementChild : option.parentNode.lastElementChild;
                        }
                        if (item.getAttribute("banswitch") != "true") {
                            options_true.querySelectorAll('SettingSetect>Options>hOption').forEach(ooption => {
                                if (ooption == option) return;
                                ooption.setAttribute('hide', '');
                            });
                            option.removeAttribute('hide');
                        }

                        item.setect = option.getAttribute("tag");
                        eval(`${item.getAttribute("onsetect").replace(/this/gi, "item")}`);
                    }
                });
            });
        } catch (error) { }
    }
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

async function GetSettings(keys) {
    try {
        const res = await fetch(`/api/settings?keys=${decodeURIComponent(JSON.stringify(keys))}`);
        if (res.status == 403) {
            throw new Error("需要登录后才可操作");
        } else if (!res.ok) throw new Error(res.statusText);
        const result = await res.json();
        return [result, null];
    } catch (error) {
        return [null, error];
    }
}

async function LoadSettings2Elements(elements) {
    try {
        if (!elements.isArray()) elements = [elements];
        let keys = [];
        for (let i = 0; i < elements.length; i++) {
            if (elements[i].isArray()) {
                let tkeys = [elements[i][0]];
                for (let i2 = 1; i2 < elements[i].length; i2++) {
                    tkeys.push(elements[i][i2].id);
                    elements[i][i2].disabled = true;
                }
                keys.push(tkeys);
            } else {
                keys.push(elements[i].id);
                elements[i].disabled = true;
            }
        }
        const [result, error] = await GetSettings(keys);
        if (error == null) {
            if (result.success) {
                for (let i = 0; i < elements.length; i++) {
                    if (elements[i].isArray()) {
                        for (let i2 = 1; i2 < elements[i].length; i2++) {
                            SetElementValue(elements[i][i2], result.data[elements[i][0]][elements[i][i2].id]);
                            elements[i][i2].disabled = false;
                        }
                    } else {
                        SetElementValue(elements[i], result.data[elements[i].id]);
                        elements[i].disabled = false;
                    }
                }
                return null;
            } else throw new Error(result.message);
        } else throw error;
    } catch (error) {
        return error;
    }
}

function SetElementValue(element, value) {
    switch (element.tagName) {
        case "INPUT":
        case "TEXTAREA":
            type = element.tagName == "TEXTAREA" ? "text" : element.type;
            switch (type) {
                case "checkbox":
                    element.checked = value;
                    break;
                default:
                    element.value = value;
                    break;
            }
            break;
        default:
            break;
    }
}

async function SaveSettings(settings) {
    try {
        const res = await fetch(`/api/settings`, {
            method: "POST",
            header: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(settings)
        })
        if (res.status == 403) {
            throw new Error("需要登录后才可操作");
        } else if (!res.ok) throw new Error(res.statusText);
        const result = await res.json();
        return [result, null];
    } catch (error) {
        return [null, error];
    }
}

async function SaveSettings4Elements(elements) {
    try {
        if (!elements.isArray()) elements = [elements];
        let keys = {};
        for (let i = 0; i < elements.length; i++) {
            if (elements[i].isArray()) {
                keys[elements[i][0]] = {};
                for (let i2 = 1; i2 < elements[i].length; i2++) {
                    keys[elements[i][0]][elements[i][i2].id] = GetElementValue(elements[i][i2]);
                    elements[i][i2].disabled = true;
                }
            } else {
                keys[elements[i].id] = GetElementValue(elements[i]);
                elements[i].disabled = true;
            }
        }
        const [result, error] = await SaveSettings(keys);
        if (error == null) {
            if (result.success) {
                for (let i = 0; i < elements.length; i++) {
                    if (elements[i].isArray()) {
                        for (let i2 = 1; i2 < elements[i].length; i2++) {
                            elements[i][i2].disabled = false;
                        }
                    } else {
                        elements[i].disabled = false;
                    }
                }
                return null;
            } else throw new Error(result.message);
        } else throw error;
    } catch (error) {
        return error;
    }
}

function GetElementValue(element) {
    switch (element.tagName) {
        case "INPUT":
        case "TEXTAREA":
            type = element.tagName == "TEXTAREA" ? "text" : element.type;
            switch (type) {
                case "checkbox":
                    return element.checked;
                default:
                    return element.value;
            }
            break;
        default:
            break;
    }
}