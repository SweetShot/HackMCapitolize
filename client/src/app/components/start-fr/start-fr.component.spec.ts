import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StartFRComponent } from './start-fr.component';

describe('StartFRComponent', () => {
  let component: StartFRComponent;
  let fixture: ComponentFixture<StartFRComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StartFRComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StartFRComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
