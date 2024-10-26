(() => {
  var __defProp = Object.defineProperty;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __moduleCache = /* @__PURE__ */ new WeakMap;
  var __toCommonJS = (from) => {
    var entry = __moduleCache.get(from), desc;
    if (entry)
      return entry;
    entry = __defProp({}, "__esModule", { value: true });
    if (from && typeof from === "object" || typeof from === "function")
      __getOwnPropNames(from).map((key) => !__hasOwnProp.call(entry, key) && __defProp(entry, key, {
        get: () => from[key],
        enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable
      }));
    __moduleCache.set(from, entry);
    return entry;
  };
  var __export = (target, all) => {
    for (var name in all)
      __defProp(target, name, {
        get: all[name],
        enumerable: true,
        configurable: true,
        set: (newValue) => all[name] = () => newValue
      });
  };

  // web/static/js/InfiniteSlider.js
  var exports_InfiniteSlider = {};
  __export(exports_InfiniteSlider, {
    InfiniteSlider: () => InfiniteSlider
  });
  function createElement(tag, prop) {
    return Object.assign(document.createElement(tag), prop);
  }
  var DEFAULT_CONFIG = {
    enableDragging: true,
    enableAutoAnimate: false,
    useMultiSlideView: false,
    skipBulletBtns: false,
    multiSlideCount: 1
  };

  class InfiniteSlider {
    containerElem;
    sliderElem;
    slidesElems;
    bulletElems = [];
    prevBtn;
    nextBtn;
    autoSlideIntvId;
    totalSlides;
    currentSlide = 0;
    isPlaying = false;
    animationDuration = 500;
    animationSleep = 1500;
    isBigSlider = false;
    isDragging = false;
    config = DEFAULT_CONFIG;
    afterAnimate;
    constructor(containerElemSelector = "[data-slider-container]", config) {
      if (config) {
        this.config = {
          ...this.config,
          ...config
        };
      }
      this.containerElem = document.querySelector(containerElemSelector);
      this.sliderElem = this.containerElem.querySelector("[data-slider-slider]");
      this.slidesElems = this.sliderElem.querySelectorAll("[data-slider-slide]");
      this.prevBtn = this.containerElem.querySelector("[data-slider-prev]");
      this.nextBtn = this.containerElem.querySelector("[data-slider-next]");
      this.animationDuration = parseInt(this.sliderElem.dataset.sliderAnimationDuration) || this.animationDuration;
      this.animationSleep = parseInt(this.sliderElem.dataset.sliderAnimationSleep) || this.animationSleep;
      this.totalSlides = this.slidesElems.length;
      this.isBigSlider = window.matchMedia("(width >= 768px)").matches;
      if (!this.config.skipBulletBtns) {
        const sliderNavigation = createElement("div", {
          className: "slider-navigation"
        });
        for (let i = 0;i < this.totalSlides; i++) {
          const btn = createElement("button", {
            type: "button",
            className: "slider-bullet",
            onclick: () => this.goto(i)
          });
          this.bulletElems.push(btn);
          sliderNavigation.append(...this.bulletElems);
          this.containerElem.append(sliderNavigation);
        }
      }
      this.multiSlideViewOn = config.useMultiSlideView;
      this.multiSlideCount = config.multiSlideCount;
      this.autoAnimate = config.enableAutoAnimate;
      this.sliderElem.addEventListener("transitionend", () => {
        if (this.currentSlide <= -1)
          this.currentSlide = this.totalSlides - 1;
        if (this.currentSlide >= this.totalSlides)
          this.currentSlide = 0;
        this.animate(0);
      });
      this.sliderElem.prepend(this.slidesElems[this.totalSlides - 1].cloneNode(true));
      this.sliderElem.append(this.slidesElems[0].cloneNode(true));
      if (this.isBigSlider&&this.multiSlideViewOn) {
        this.sliderElem.style.transform = `translateX(133.33%)`;
        this.sliderElem.prepend(this.slidesElems[this.totalSlides - 2].cloneNode(true));
        this.sliderElem.append(this.slidesElems[1].cloneNode(true));
      }
      this.animate = this.animate.bind(this);
      this.prev = this.prev.bind(this);
      this.next = this.next.bind(this);
      this.goto = this.goto.bind(this);
      this.play = this.play.bind(this);
      this.prevBtn.addEventListener("click", (e) => {
        e.preventDefault();
        this.stop();
        this.prev();
        if (this.config.enableAutoAnimate) {
          this.play();
        }
      });
      this.nextBtn.addEventListener("click", (e) => {
        e.preventDefault();
        this.stop();
        this.next();
        if (this.config.enableAutoAnimate) {
          this.play();
        }
      });
      if (this.config.enableDragging) {
        let startX = 0;
        let accDrag = 0;
        let slideThreshold = Math.floor(this.slidesElems[0].scrollWidth / 2);
        this.containerElem.addEventListener("pointerdown", (e) => {
          this.containerElem.style.cursor = "grabbing";
          this.isDragging = true;
          startX = e.x;
        });
        this.containerElem.addEventListener("pointermove", (e) => {
          if (this.isDragging) {
            accDrag += e.movementX;
            const translateX = this.multiSlideViewOn ? this.calcTranslateAmount(this.isBigSlider) : 100;
            this.sliderElem.style.transform = `translateX(calc(${(-this.currentSlide - 1) * translateX}% + ${accDrag}px))`;
            if (Math.abs(accDrag) >= slideThreshold) {
              return;
            }
          }
        });
        this.containerElem.addEventListener("pointerup", () => {
          this.containerElem.style.cursor = "auto";
          if (Math.abs(accDrag) >= slideThreshold) {
            if (0 < accDrag) {
              this.prev();
            } else {
              this.next();
            }
          } else {
            this.animate();
          }
          startX = 0;
          accDrag = 0;
          this.isDragging = false;
        });
      }
      if (this.config.enableAutoAnimate) {
        this.containerElem.addEventListener("pointerenter", () => {
          this.stop();
        });
        this.containerElem.addEventListener("pointerleave", () => {
          this.play();
        });
        this.play();
      }
    }
    animate(ms = this.animationDuration) {
      const cMod = this.calcSlideMod(this.currentSlide, this.totalSlides);
      this.sliderElem.style.transitionDuration = `${ms}ms`;
      this.translateSlideToCenter();
      this.slidesElems.forEach((slide, i) => slide.classList.toggle("is-active", cMod === i));
      this.bulletElems.forEach((bulletBtn, i) => bulletBtn.classList.toggle("is-active", cMod === i));
      if (typeof this.afterAnimate === "function") {
        this.afterAnimate(this);
      }
    }
    prev() {
      if (this.currentSlide <= -1)
        return;
      this.currentSlide -= 1;
      this.animate();
    }
    next() {
      if (this.currentSlide >= this.totalSlides)
        return;
      this.currentSlide += 1;
      this.animate();
    }
    goto(index) {
      this.currentSlide = index;
      this.animate();
    }
    play() {
      if (this.autoSlideIntvId) {
        clearInterval(this.autoSlideIntvId);
      }
      this.autoSlideIntvId = setInterval(this.next, this.animationSleep + this.animationDuration);
    }
    stop() {
      clearInterval(this.autoSlideIntvId);
      this.autoSlideIntvId = undefined;
    }
    translateSlideToCenter() {
      const translateAmount = this.multiSlideViewOn ? this.calcTranslateAmount(this.isBigSlider) : 100;
      this.sliderElem.style.transform = `translateX(${(-this.currentSlide - 1) * translateAmount}%)`;
    }
    calcTranslateAmount(isBigSlider) {
      if (isBigSlider) {
        return 33.333;
      } else {
        return 100;
      }
    }
    calcSlideMod(currentIndex, totalSlides) {
      return (currentIndex % totalSlides + totalSlides) % totalSlides;
    }
  }
  if (typeof window !== "undefined") {
    window["InfiniteSlider"] = InfiniteSlider;
  }
})();
